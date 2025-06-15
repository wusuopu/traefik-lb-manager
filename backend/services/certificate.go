package services

import (
	"app/config"
	"app/di"
	"app/models"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/gogf/gf/v2/util/gconv"
)

const system_setting_account_key = "acme_account"

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}
func (u *AcmeUser) GetEmail() string {
	return u.Email
}
func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func NewAcmeUser(email string) (*AcmeUser, error) {
  user := &AcmeUser{}
	var err error = nil

	// 从数据库中读取账号信息
	var settings models.SystemSetting
	results := di.Container.DB.Where("name = ?", system_setting_account_key).First(&settings)
	if results.Error == nil {
		data := gconv.MapDeep(settings.Value)
		if data != nil {
			gconv.Struct(data, &user)
		}
	}

	// 创建新账号
	if user.Email == "" || user.Key == "" {
		user.Email = email
		user.Key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			di.Container.Logger.Error(fmt.Sprintf("generate User PrivateKey failed: %s", err))
			return nil, err
		}
	}

	return user, nil
}

// ====================================================================================
type DBProvider struct {
	cert		*models.Certificate
}
func (p *DBProvider) Present(domain, token, keyAuth string) error {
	if p.cert == nil {
		return fmt.Errorf("certificate is nil")
	}

	// 查找最新一条待处理的数据
	var obj models.Certificate
	results := di.Container.DB.First(&obj, p.cert.ID)

	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("DBProvider Present error: %s %s", domain, results.Error))
		return results.Error
	}
	if domain != obj.Domain {
		return fmt.Errorf("domain %s not match %s", domain, obj.Domain)
	}

	results = di.Container.DB.
		Model(&models.Certificate{}).
		Where("id = ?", obj.ID).
		Updates(map[string]interface{}{
			"acme_token": token,
			"acme_key_auth": keyAuth,
		})
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("DBProvider Present token error: %s %s", domain, results.Error))
		return results.Error
	}

	return nil
}
func (p *DBProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}

// ====================================================================================
/*
	生成证书步聚：
	1. 创建用户
		user, err := NewAcmeUser(email)
	2. 创建证书管理器
		certManager, err := NewAcmeCertManager(user, "")
	3. 注册账号
		err := certManager.Register()
	4. 生成证书
		err := certManager.ObtainCertificate(cert)
*/
type AcmeCertManager struct {
	user				*AcmeUser
	client			*lego.Client
	acmeServer	string
}
func NewAcmeCertManager(user *AcmeUser, acmeServer string) (*AcmeCertManager, error) {
	cfg := lego.NewConfig(user)
	if acmeServer != "" {
		cfg.CADirURL = acmeServer
	} else {
		if gin.Mode() == gin.ReleaseMode {
			cfg.CADirURL = lego.LEDirectoryProduction
		} else {
			cfg.CADirURL = lego.LEDirectoryStaging
		}
	}
	cfg.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(cfg)
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("create acme client failed: %s", err))
		return nil, err
	}

	// 设置 HTTP-01 challenge
	err = client.Challenge.SetHTTP01Provider(&DBProvider{})
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("SetHTTP01Provider error: %s", err))
		return nil, err
	}

	return &AcmeCertManager{
		user: user,
		client: client,
		acmeServer: acmeServer,
	}, nil
}

func (m *AcmeCertManager) Register() (error) {
	// 用户已注册
	var reg *registration.Resource
	var err error

	if m.user.Registration != nil && m.user.Registration.URI != "" {
		reg, err = m.client.Registration.QueryRegistration()
	}
	if reg == nil {
    di.Container.Logger.Info("Acme Register New Account")
		reg, err = m.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	}

  if err != nil {
    di.Container.Logger.Error(fmt.Sprintf("register acme user failed: %s", err))
    return err
  }

  m.user.Registration = reg

	// 账号信息保存到数据库
	var settings models.SystemSetting
	di.Container.DB.Where("name = ?", system_setting_account_key).First(&settings)

	settings.Name = system_setting_account_key
	settings.Value, _ = json.Marshal(gconv.MapDeep(m.user))
	if settings.ID > 0 {
		di.Container.DB.Save(&settings)
	} else {
		di.Container.DB.Create(&settings)
	}
	
  return nil
}

func (m *AcmeCertManager) ObtainCertificate(cert *models.Certificate) error {
	if cert == nil || cert.ID == 0 {
		return errors.New("certificate is nil")
	}
	err := m.client.Challenge.SetHTTP01Provider(&DBProvider{cert: cert})
	if err != nil {
		return err
	}

	// 请求证书
	request := certificate.ObtainRequest{
		Domains: []string{cert.Domain},
		Bundle:  true,
	}

	results := di.Container.DB.
		Model(&models.Certificate{}).
		Where("id = ?", cert.ID).
		Updates(map[string]interface{}{
			"status": models.CERTIFICATE_STATUS_PENDING,
		})
	if results.Error != nil {
		return results.Error
	}

	di.Container.Logger.Info(fmt.Sprintf("start to obtain certificate for %s", cert.Domain))
	certificates, err := m.client.Certificate.Obtain(request)
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("obtain certificate failed for %s %s", cert.Domain, err))
		// 更新状态
		di.Container.DB.
			Model(&models.Certificate{}).
			Where("id = ?", cert.ID).
			Updates(map[string]interface{}{
				"status": models.CERTIFICATE_STATUS_FAILED,
			})
		return err
	}

	payload := map[string]interface{}{
		"cert": string(certificates.Certificate),
		"key":   string(certificates.PrivateKey),
		"status": models.CERTIFICATE_STATUS_COMPLETE,
	}

	// 解析证书有效期
	certInfo, _ := m.ParseCertificate(certificates.Certificate)
	if certInfo != nil {
		payload["effective_at"] = certInfo.NotBefore
		payload["expired_at"] = certInfo.NotAfter
	}

	// 保存证书
	results = di.Container.DB.
		Model(&models.Certificate{}).
		Where("id = ?", cert.ID).
		Updates(payload)
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("save certificate failed for %s %s", cert.Domain, results.Error))
		return results.Error
	}

	return nil
}

func (m *AcmeCertManager) ParseCertificate(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, errors.New("invalid certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse certificate failed: %s", err)
	}

  return cert, nil
}

// ====================================================================================

type CertificateService struct {
}
func (c *CertificateService) Obtain(cert *models.Certificate) error {
	user, err := NewAcmeUser(config.Config.SSL.Email)
	if err != nil {
		return err
	}

	certManager, err := NewAcmeCertManager(user, "")
	if err != nil {
		return err
	}

	err = certManager.Register()
	if err != nil {
		return err
	}

	err = certManager.ObtainCertificate(cert)

	return err
}