import _ from 'lodash';

export const allMiddlewares = [
  "addPrefix",
  "basicAuth",
  "buffering",
  "chain",
  "circuitBreaker",
  "compress",
  "contentType",
  "digestAuth",
  "errors",
  "forwardAuth",
  "grpcWeb",
  "headers",
  "ipWhiteList",
  "ipAllowList",
  "inFlightReq",
  "passTLSClientCert",
  "plugin",
  "rateLimit",
  "redirectRegex",
  "redirectScheme",
  "replacePath",
  "replacePathRegex",
  "retry",
  "stripPrefix",
  "stripPrefixRegex"
]

function addPropertiesTitle (properties: any) {
  for (const key in properties) {
    const item = properties[key]
    if (!item.title) { item.title = key }

    if (item.properties) {
      addPropertiesTitle(item.properties)
    }
  }
}

// https://json.schemastore.org/traefik-v3-file-provider.json
const allSchemas: { [key: string]: any} = {
  "addPrefixMiddleware": {
    "type": "object",
    "description": "The AddPrefix middleware updates the URL Path of the request before forwarding it.",
    "properties": {
      "prefix": {
        "type": "string",
        "description": "prefix is the string to add before the current path in the requested URL. It should include the leading slash (/)."
      }
    },
    "additionalProperties": false
  },
  "basicAuthMiddleware": {
    "type": "object",
    "description": "The BasicAuth middleware is a quick way to restrict access to your services to known users. If both users and usersFile are provided, the two are merged. The contents of usersFile have precedence over the values in users.",
    "properties": {
      "users": {
        "type": "array",
        "description": "The users option is an array of authorized users. Each user will be declared using the `name:hashed-password` format.",
        "items": {
          "type": "string"
        }
      },
      "usersFile": {
        "type": "string",
        "description": "The usersFile option is the path to an external file that contains the authorized users for the middleware.\n\nThe file content is a list of `name:hashed-password`."
      },
      "realm": {
        "type": "string",
        "description": "You can customize the realm for the authentication with the realm option. The default value is traefik.",
        "default": "traefik"
      },
      "headerField": {
        "type": "string",
        "description": "You can define a header field to store the authenticated user using the headerField option."
      },
      "removeHeader": {
        "type": "boolean",
        "description": "Set the removeHeader option to true to remove the authorization header before forwarding the request to your service. (Default value is false.)",
        "default": false
      }
    }
  },
  "bufferingMiddleware": {
    "type": "object",
    "description": "The Buffering middleware gives you control on how you want to read the requests before sending them to services.\n\nWith Buffering, Traefik reads the entire request into memory (possibly buffering large requests into disk), and rejects requests that are over a specified limit.\n\nThis can help services deal with large data (multipart/form-data for example), and can minimize time spent sending data to a service.",
    "properties": {
      "maxRequestBodyBytes": {
        "type": "integer",
        "description": "With the maxRequestBodyBytes option, you can configure the maximum allowed body size for the request (in Bytes).\n\nIf the request exceeds the allowed size, it is not forwarded to the service and the client gets a 413 (Request Entity Too Large) response."
      },
      "memRequestBodyBytes": {
        "type": "integer",
        "description": "You can configure a threshold (in Bytes) from which the request will be buffered on disk instead of in memory with the memRequestBodyBytes option."
      },
      "maxResponseBodyBytes": {
        "type": "integer",
        "description": "With the maxResponseBodyBytes option, you can configure the maximum allowed response size from the service (in Bytes).\n\nIf the response exceeds the allowed size, it is not forwarded to the client. The client gets a 413 (Request Entity Too Large) response instead."
      },
      "memResponseBodyBytes": {
        "type": "integer",
        "description": "You can configure a threshold (in Bytes) from which the response will be buffered on disk instead of in memory with the memResponseBodyBytes option."
      },
      "retryExpression": {
        "type": "string",
        "description": "You can have the Buffering middleware replay the request with the help of the retryExpression option."
      }
    },
    "additionalProperties": false
  },
  "chainMiddleware": {
    "type": "object",
    "description": "The Chain middleware enables you to define reusable combinations of other pieces of middleware. It makes reusing the same groups easier.",
    "properties": {
      "middlewares": {
        "type": "array",
        "minItems": 1,
        "items": {
          "type": "string"
        }
      }
    },
    "additionalProperties": false
  },
  "circuitBreakerMiddleware": {
    "type": "object",
    "description": "The circuit breaker protects your system from stacking requests to unhealthy services (resulting in cascading failures).\n\nWhen your system is healthy, the circuit is closed (normal operations). When your system becomes unhealthy, the circuit becomes open and the requests are no longer forwarded (but handled by a fallback mechanism).\n\nTo assess if your system is healthy, the circuit breaker constantly monitors the services.",
    "properties": {
      "expression": {
        "type": "string",
        "description": "You can specify an expression that, once matched, will trigger the circuit breaker (and apply the fallback mechanism instead of calling your services)."
      },
      "checkPeriod": {
        "type": "string",
        "description": "The interval between successive checks of the circuit breaker condition (when in standby state)"
      },
      "fallbackDuration": {
        "type": "string",
        "description": "The duration for which the circuit breaker will wait before trying to recover (from a tripped state)."
      },
      "recoveryDuration": {
        "type": "string",
        "description": "The duration for which the circuit breaker will try to recover (as soon as it is in recovering state)."
      },
      "responseCode": {
        "type": "integer",
        "description": "The status code that the circuit breaker will return while it is in the open state."
      }
    },
    "additionalProperties": false
  },
  "compressMiddleware": {
    "type": "object",
    "description": "The Compress middleware enables the gzip compression.",
    "properties": {
      "excludedContentTypes": {
        "type": "array",
        "description": "excludedContentTypes specifies a list of content types to compare the Content-Type header of the incoming requests to before compressing.\n\nThe requests with content types defined in excludedContentTypes are not compressed.\n\nContent types are compared in a case-insensitive, whitespace-ignored manner.",
        "items": {
          "type": "string"
        }
      },
      "minResponseBodyBytes": {
        "description": "specifies the minimum amount of bytes a response body must have to be compressed.",
        "type": "integer"
      },
      "defaultEncoding": {
        "type": "string",
        "description": "defaultEncoding specifies the default encoding if the Accept-Encoding header is not in the request or contains a wildcard (*)."
      },
      "includedContentTypes": {
        "type": "array",
        "description": "includedContentTypes specifies a list of content types to compare the Content-Type header of the responses before compressing.\n\nThe responses with content types defined in includedContentTypes are compressed.\n\nContent types are compared in a case-insensitive, whitespace-ignored manner.",
        "items": {
          "type": "string"
        }
      },
      "encodings": {
        "type": "array",
        "description": "encodings specifies the list of supported compression encodings. At least one encoding value must be specified, and valid entries are zstd (Zstandard), br (Brotli), and gzip (Gzip). The order of the list also sets the priority, the top entry has the highest priority.",
        "items": {
          "type": "string"
        }
      }
    },
    "additionalProperties": false
  },
  "contentTypeMiddleware": {
    "type": "object",
    "description": "The Content-Type middleware - or rather its unique autoDetect option - specifies whether to let the Content-Type header, if it has not been set by the backend, be automatically set to a value derived from the contents of the response.\n\nAs a proxy, the default behavior should be to leave the header alone, regardless of what the backend did with it. However, the historic default was to always auto-detect and set the header if it was nil, and it is going to be kept that way in order to support users currently relying on it. This middleware exists to enable the correct behavior until at least the default one can be changed in a future version.",
    "properties": {
      "autoDetect": {
        "type": "boolean",
        "description": "autoDetect specifies whether to let the Content-Type header, if it has not been set by the backend, be automatically set to a value derived from the contents of the response.",
        "default": false
      }
    },
    "additionalProperties": false
  },
  "digestAuthMiddleware": {
    "type": "object",
    "description": "The DigestAuth middleware is a quick way to restrict access to your services to known users. If both users and usersFile are provided, the two are merged. The contents of usersFile have precedence over the values in users.",
    "properties": {
      "users": {
        "type": "array",
        "description": "The users option is an array of authorized users. Each user will be declared using the `name:realm:encoded-password` format.",
        "items": {
          "type": "string"
        }
      },
      "usersFile": {
        "type": "string",
        "description": "The usersFile option is the path to an external file that contains the authorized users for the middleware.\n\nThe file content is a list of `name:realm:encoded-password`."
      },
      "realm": {
        "type": "string",
        "description": "You can customize the realm for the authentication with the realm option. The default value is traefik.",
        "default": "traefik"
      },
      "headerField": {
        "type": "string",
        "description": "You can customize the header field for the authenticated user using the headerField option."
      },
      "removeHeader": {
        "type": "boolean",
        "description": "Set the removeHeader option to true to remove the authorization header before forwarding the request to your service. (Default value is false.)",
        "default": false
      }
    },
    "additionalProperties": false
  },
  "errorsMiddleware": {
    "type": "object",
    "description": "The ErrorPage middleware returns a custom page in lieu of the default, according to configured ranges of HTTP Status codes. The error page itself is not hosted by Traefik.",
    "properties": {
      "status": {
        "type": "array",
        "description": "The status that will trigger the error page.\n\nThe status code ranges are inclusive (500-599 will trigger with every code between 500 and 599, 500 and 599 included). You can define either a status code like 500 or ranges with a syntax like 500-599.",
        "items": {
          "type": "string"
        }
      },
      "service": {
        "type": "string",
        "description": "The service that will serve the new requested error page."
      },
      "query": {
        "type": "string",
        "description": "The URL for the error page (hosted by service). You can use {status} in the query, that will be replaced by the received status code."
      }
    },
    "additionalProperties": false
  },
  "forwardAuthMiddleware": {
    "type": "object",
    "description": "The ForwardAuth middleware delegate the authentication to an external service. If the service response code is 2XX, access is granted and the original request is performed. Otherwise, the response from the authentication server is returned.",
    "properties": {
      "address": {
        "type": "string",
        "description": "The address option defines the authentication server address."
      },
      "tls": {
        "type": "object",
        "description": "The tls option is the TLS configuration from Traefik to the authentication server.",
        "properties": {
          "ca": {
            "type": "string",
            "description": "Certificate Authority used for the secured connection to the authentication server."
          },
          "cert": {
            "type": "string",
            "description": "Public certificate used for the secured connection to the authentication server."
          },
          "key": {
            "type": "string",
            "description": "Private certificate used for the secure connection to the authentication server."
          },
          "insecureSkipVerify": {
            "type": "boolean",
            "description": "If insecureSkipVerify is true, TLS for the connection to authentication server accepts any certificate presented by the server and any host name in that certificate."
          }
        }
      },
      "trustForwardHeader": {
        "type": "boolean",
        "description": "Set the trustForwardHeader option to true to trust all the existing X-Forwarded-* headers."
      },
      "authResponseHeaders": {
        "type": "array",
        "description": "The authResponseHeaders option is the list of the headers to copy from the authentication server to the request.",
        "items": {
          "type": "string"
        }
      },
      "authResponseHeadersRegex": {
        "type": "string",
        "description": "The authResponseHeadersRegex option is the regex to match headers to copy from the authentication server response and set on forwarded request, after stripping all headers that match the regex."
      },
      "authRequestHeaders": {
        "type": "array",
        "description": "The authRequestHeaders option is the list of the headers to copy from the request to the authentication server.",
        "items": {
          "type": "string"
        }
      },
      "addAuthCookiesToResponse": {
        "type": "array",
        "description": "The addAuthCookiesToResponse option is the list of cookies to copy from the authentication server to the response, replacing any existing conflicting cookie from the forwarded response.",
        "items": {
          "type": "string"
        }
      }
    },
    "additionalProperties": false
  },
  "grpcWebMiddleware": {
    "type": "object",
    "description": "The GrpcWeb middleware converts gRPC Web requests to HTTP/2 gRPC requests before forwarding them to the backends.",
    "properties": {
      "allowOrigins": {
        "type": "array",
        "description": "The allowOrigins contains the list of allowed origins. A wildcard origin * can also be configured to match all requests.",
        "items": {
          "type": "string"
        }
      }
    },
    "additionalProperties": false
  },
  "headersMiddleware": {
    "type": "object",
    "description": "The Headers middleware can manage the requests/responses headers.",
    "properties": {
      "customRequestHeaders": {
        "type": "object",
        "description": "The customRequestHeaders option lists the Header names and values to apply to the request.",
        "additionalProperties": {
          "type": "string"
        }
      },
      "customResponseHeaders": {
        "type": "object",
        "description": "The customResponseHeaders option lists the Header names and values to apply to the response.",
        "additionalProperties": {
          "type": "string"
        }
      },
      "accessControlAllowCredentials": {
        "type": "boolean",
        "description": "The accessControlAllowCredentials indicates whether the request can include user credentials."
      },
      "accessControlAllowHeaders": {
        "type": "array",
        "description": "The accessControlAllowHeaders indicates which header field names can be used as part of the request.",
        "items": {
          "type": "string"
        }
      },
      "accessControlAllowMethods": {
        "type": "array",
        "description": "The accessControlAllowMethods indicates which methods can be used during requests.",
        "items": {
          "type": "string"
        }
      },
      "accessControlAllowOriginList": {
        "type": "array",
        "description": "The accessControlAllowOriginList indicates whether a resource can be shared by returning different values.\n\nA wildcard origin * can also be configured, and will match all requests. If this value is set by a backend server, it will be overwritten by Traefik\n\nThis value can contain a list of allowed origins.",
        "items": {
          "type": "string"
        }
      },
      "accessControlAllowOriginListRegex": {
        "type": "array",
        "description": "The accessControlAllowOriginListRegex option is the counterpart of the accessControlAllowOriginList option with regular expressions instead of origin values.",
        "items": {
          "type": "string"
        }
      },
      "accessControlExposeHeaders": {
        "type": "array",
        "description": "The accessControlExposeHeaders indicates which headers are safe to expose to the api of a CORS API specification.",
        "items": {
          "type": "string"
        }
      },
      "accessControlMaxAge": {
        "type": "integer",
        "description": "The accessControlMaxAge indicates how long (in seconds) a preflight request can be cached."
      },
      "addVaryHeader": {
        "type": "boolean",
        "description": "The addVaryHeader is used in conjunction with accessControlAllowOriginList to determine whether the vary header should be added or modified to demonstrate that server responses can differ based on the value of the origin header."
      },
      "allowedHosts": {
        "type": "array",
        "description": "The allowedHosts option lists fully qualified domain names that are allowed.",
        "items": {
          "type": "string"
        }
      },
      "hostsProxyHeaders": {
        "type": "array",
        "description": "The hostsProxyHeaders option is a set of header keys that may hold a proxied hostname value for the request.",
        "items": {
          "type": "string"
        }
      },
      "sslRedirect": {
        "type": "boolean",
        "description": "The sslRedirect is set to true, then only allow https requests."
      },
      "sslTemporaryRedirect": {
        "type": "boolean",
        "description": "Set the sslTemporaryRedirect to true to force an SSL redirection using a 302 (instead of a 301)."
      },
      "sslHost": {
        "type": "string",
        "description": "The sslHost option is the host name that is used to redirect http requests to https."
      },
      "sslProxyHeaders": {
        "type": "object",
        "description": "The sslProxyHeaders option is set of header keys with associated values that would indicate a valid https request. Useful when using other proxies with header like: \"X-Forwarded-Proto\": \"https\".",
        "additionalProperties": {
          "type": "string"
        }
      },
      "sslForceHost": {
        "type": "boolean",
        "description": "Set sslForceHost to true and set SSLHost to forced requests to use SSLHost even the ones that are already using SSL."
      },
      "stsSeconds": {
        "type": "integer",
        "description": "The stsSeconds is the max-age of the Strict-Transport-Security header. If set to 0, would NOT include the header."
      },
      "stsIncludeSubdomains": {
        "type": "boolean",
        "description": "The stsIncludeSubdomains is set to true, the includeSubDomains directive will be appended to the Strict-Transport-Security header."
      },
      "stsPreload": {
        "type": "boolean",
        "description": "Set stsPreload to true to have the preload flag appended to the Strict-Transport-Security header."
      },
      "forceSTSHeader": {
        "type": "boolean",
        "description": "Set forceSTSHeader to true, to add the STS header even when the connection is HTTP."
      },
      "frameDeny": {
        "type": "boolean",
        "description": "Set frameDeny to true to add the X-Frame-Options header with the value of DENY."
      },
      "customFrameOptionsValue": {
        "type": "string",
        "description": "The customFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option."
      },
      "contentTypeNosniff": {
        "type": "boolean",
        "description": "Set contentTypeNosniff to true to add the X-Content-Type-Options header with the value nosniff."
      },
      "browserXssFilter": {
        "type": "boolean",
        "description": "Set browserXssFilter to true to add the X-XSS-Protection header with the value 1; mode=block."
      },
      "customBrowserXSSValue": {
        "type": "string",
        "description": "The customBrowserXssValue option allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option."
      },
      "contentSecurityPolicy": {
        "type": "string",
        "description": "The contentSecurityPolicy option allows the Content-Security-Policy header value to be set with a custom value."
      },
      "contentSecurityPolicyReportOnly": {
        "type": "string",
        "description": "The contentSecurityPolicyReportOnly option allows the Content-Security-Policy-Report-Only header value to be set with a custom value."
      },
      "publicKey": {
        "type": "string",
        "description": "The publicKey implements HPKP to prevent MITM attacks with forged certificates."
      },
      "referrerPolicy": {
        "type": "string",
        "description": "The referrerPolicy allows sites to control when browsers will pass the Referer header to other sites."
      },
      "featurePolicy": {
        "type": "string",
        "description": "The featurePolicy allows sites to control browser features."
      },
      "permissionsPolicy": {
        "type": "string",
        "description": "The permissionsPolicy allows sites to control browser features."
      },
      "isDevelopment": {
        "type": "boolean",
        "description": "Set isDevelopment to true when developing. The AllowedHosts, SSL, and STS options can cause some unwanted effects. Usually testing happens on http, not https, and on localhost, not your production domain.\nIf you would like your development environment to mimic production with complete Host blocking, SSL redirects, and STS headers, leave this as false."
      }
    },
    "additionalProperties": false
  },
  "ipStrategy": {
    "type": "object",
    "description": "The ipStrategy option defines parameters that set how Traefik will determine the client IP.",
    "properties": {
      "depth": {
        "type": "integer",
        "description": "The depth option tells Traefik to use the X-Forwarded-For header and take the IP located at the depth position (starting from the right). If depth is greater than the total number of IPs in X-Forwarded-For, then the client IP will be empty. depth is ignored if its value is lesser than or equal to 0."
      },
      "excludedIPs": {
        "type": "array",
        "description": "excludedIPs tells Traefik to scan the X-Forwarded-For header and pick the first IP not in the list. If depth is specified, excludedIPs is ignored.",
        "items": {
          "type": "string"
        }
      }
    },
    "additionalProperties": false
  },
  "ipWhiteListMiddleware": {
    "type": "object",
    "description": "DEPRECATED: IPWhitelist accepts / refuses requests based on the client IP.",
    "properties": {
      "sourceRange": {
        "type": "array",
        "description": "The sourceRange option sets the allowed IPs (or ranges of allowed IPs by using CIDR notation).",
        "items": {
          "type": "string"
        }
      },
      "ipStrategy": {
        "$ref": "#/definitions/ipStrategy"
      }
    },
    "additionalProperties": false
  },
  "ipAllowListMiddleware": {
    "type": "object",
    "description": "IPAllowList accepts / refuses requests based on the client IP.",
    "properties": {
      "sourceRange": {
        "type": "array",
        "description": "The sourceRange option sets the allowed IPs (or ranges of allowed IPs by using CIDR notation).",
        "items": {
          "type": "string"
        }
      },
      "rejectStatusCode": {
        "type": "integer",
        "description": "RejectStatusCode defines the HTTP status code used for refused requests. If not set, the default is 403 (Forbidden)."
      },
      "ipStrategy": {
        "$ref": "#/definitions/ipStrategy"
      }
    },
    "additionalProperties": false
  },
  "sourceCriterion": {
    "type": "object",
    "description": "SourceCriterion defines what criterion is used to group requests as originating from a common source. The precedence order is ipStrategy, then requestHeaderName, then requestHost. If none are set, the default is to use the requestHost.",
    "properties": {
      "ipStrategy": {
        "$ref": "#/definitions/ipStrategy"
      },
      "requestHeaderName": {
        "type": "string",
        "description": "Requests having the same value for the given header are grouped as coming from the same source."
      },
      "requestHost": {
        "type": "boolean",
        "description": "Whether to consider the request host as the source."
      }
    },
    "additionalProperties": false
  },
  "inFlightReqMiddleware": {
    "type": "object",
    "description": "To proactively prevent services from being overwhelmed with high load, a limit on the number of simultaneous in-flight requests can be applied.",
    "properties": {
      "amount": {
        "type": "integer",
        "description": "The amount option defines the maximum amount of allowed simultaneous in-flight request. The middleware will return an HTTP 429 Too Many Requests if there are already amount requests in progress (based on the same sourceCriterion strategy)."
      },
      "sourceCriterion": {
        "$ref": "#/definitions/sourceCriterion"
      }
    },
    "additionalProperties": false
  },
  "passTLSClientCertMiddleware": {
    "type": "object",
    "description": "PassTLSClientCert adds in header the selected data from the passed client tls certificate.",
    "properties": {
      "pem": {
        "type": "boolean",
        "description": "The pem option sets the X-Forwarded-Tls-Client-Cert header with the escape certificate."
      },
      "info": {
        "type": "object",
        "description": "The info option select the specific client certificate details you want to add to the X-Forwarded-Tls-Client-Cert-Info header. The value of the header will be an escaped concatenation of all the selected certificate details.",
        "properties": {
          "notAfter": {
            "type": "boolean",
            "description": "Set the notAfter option to true to add the Not After information from the Validity part."
          },
          "notBefore": {
            "type": "boolean",
            "description": "Set the notBefore option to true to add the Not Before information from the Validity part."
          },
          "sans": {
            "type": "boolean",
            "description": "Set the sans option to true to add the Subject Alternative Name information from the Subject Alternative Name part."
          },
          "subject": {
            "type": "object",
            "description": "The subject select the specific client certificate subject details you want to add to the X-Forwarded-Tls-Client-Cert-Info header.",
            "properties": {
              "country": {
                "type": "boolean",
                "description": "Set the country option to true to add the country information into the subject."
              },
              "province": {
                "type": "boolean",
                "description": "Set the province option to true to add the province information into the subject."
              },
              "locality": {
                "type": "boolean",
                "description": "Set the locality option to true to add the locality information into the subject."
              },
              "organization": {
                "type": "boolean",
                "description": "Set the organization option to true to add the organization information into the subject."
              },
              "commonName": {
                "type": "boolean",
                "description": "Set the commonName option to true to add the commonName information into the subject."
              },
              "serialNumber": {
                "type": "boolean",
                "description": "Set the serialNumber option to true to add the serialNumber information into the subject."
              },
              "domainComponent": {
                "type": "boolean",
                "description": "Set the domainComponent option to true to add the domainComponent information into the subject."
              }
            }
          },
          "issuer": {
            "type": "object",
            "description": "The issuer select the specific client certificate issuer details you want to add to the X-Forwarded-Tls-Client-Cert-Info header.",
            "properties": {
              "country": {
                "type": "boolean",
                "description": "Set the country option to true to add the country information into the issuer."
              },
              "province": {
                "type": "boolean",
                "description": "Set the province option to true to add the province information into the issuer."
              },
              "locality": {
                "type": "boolean",
                "description": "Set the locality option to true to add the locality information into the issuer."
              },
              "organization": {
                "type": "boolean",
                "description": "Set the organization option to true to add the organization information into the issuer."
              },
              "commonName": {
                "type": "boolean",
                "description": "Set the commonName option to true to add the commonName information into the issuer."
              },
              "serialNumber": {
                "type": "boolean",
                "description": "Set the serialNumber option to true to add the serialNumber information into the issuer."
              },
              "domainComponent": {
                "type": "boolean",
                "description": "Set the domainComponent option to true to add the domainComponent information into the issuer."
              }
            }
          }
        }
      }
    },
    "additionalProperties": false
  },
  "pluginMiddleware": {
    "type": "object",
    "description": "Some plugins will need to be configured by adding a dynamic configuration.",
    "additionalProperties": {
      "type": "object"
    }
  },
  "rateLimitMiddleware": {
    "type": "object",
    "description": "The RateLimit middleware ensures that services will receive a fair number of requests, and allows one to define what fair is.",
    "properties": {
      "average": {
        "description": "average is the maximum rate, by default in requests by second, allowed for the given source.\n\nIt defaults to 0, which means no rate limiting.\n\nThe rate is actually defined by dividing average by period. So for a rate below 1 req/s, one needs to define a period larger than a second.",
        "oneOf": [
          {
            "type": "string"
          },
          {
            "type": "number"
          }
        ]
      },
      "period": {
        "oneOf": [
          {
            "type": "string"
          },
          {
            "type": "number",
            "default": 1
          }
        ],
        "description": "period, in combination with average, defines the actual maximum rate.\n\nIt defaults to 1 second."
      },
      "burst": {
        "type": "number",
        "description": "burst is the maximum number of requests allowed to go through in the same arbitrarily small period of time.\n\nIt defaults to 1.",
        "default": 1
      },
      "sourceCriterion": {
        "$ref": "#/definitions/sourceCriterion"
      }
    },
    "additionalProperties": false
  },
  "redirectRegexMiddleware": {
    "type": "object",
    "description": "RegexRedirect redirect a request from an url to another with regex matching and replacement.",
    "properties": {
      "permanent": {
        "type": "boolean",
        "description": "Set the permanent option to true to apply a permanent redirection."
      },
      "regex": {
        "type": "string",
        "description": "The regex option is the regular expression to match and capture elements from the request URL."
      },
      "replacement": {
        "type": "string",
        "description": "The replacement option defines how to modify the URL to have the new target URL. Care should be taken when defining replacement expand variables: $1x is equivalent to ${1x}, not ${1}x (see Regexp.Expand), so use ${1} syntax."
      }
    },
    "additionalProperties": false
  },
  "redirectSchemeMiddleware": {
    "type": "object",
    "description": "RedirectScheme redirect request from a scheme to another.",
    "properties": {
      "permanent": {
        "type": "boolean",
        "description": "Set the permanent option to true to apply a permanent redirection."
      },
      "scheme": {
        "type": "string",
        "description": "The scheme option defines the scheme of the new url."
      },
      "port": {
        "type": "string",
        "description": "The port option defines the port of the new url. Port in this configuration is a string, not a numeric value."
      }
    },
    "additionalProperties": false
  },
  "replacePathMiddleware": {
    "type": "object",
    "description": "Replace the path of the request url. It will replace the actual path by the specified one and will store the original path in a X-Replaced-Path header.",
    "properties": {
      "path": {
        "type": "string",
        "description": "The path option defines the path to use as replacement in the request url."
      }
    },
    "additionalProperties": false
  },
  "replacePathRegexMiddleware": {
    "type": "object",
    "description": "The ReplaceRegex replace a path from an url to another with regex matching and replacement. It will replace the actual path by the specified one and store the original path in a X-Replaced-Path header.",
    "properties": {
      "regex": {
        "type": "string",
        "description": "The regex option is the regular expression to match and capture the path from the request URL."
      },
      "replacement": {
        "type": "string",
        "description": "The replacement option defines how to modify the path to have the new target path. Care should be taken when defining replacement expand variables: $1x is equivalent to ${1x}, not ${1}x (see Regexp.Expand), so use ${1} syntax."
      }
    },
    "additionalProperties": false
  },
  "retryMiddleware": {
    "type": "object",
    "description": "The Retry middleware is in charge of reissuing a request a given number of times to a backend server if that server does not reply. To be clear, as soon as the server answers, the middleware stops retrying, regardless of the response status.",
    "properties": {
      "attempts": {
        "type": "integer",
        "description": "The attempts option defines how many times the request should be retried."
      },
      "initialInterval": {
        "type": "string",
        "description": "The initialInterval option defines the first wait time in the exponential backoff series."
      }
    },
    "additionalProperties": false,
    "required": [
      "attempts"
    ]
  },
  "stripPrefixMiddleware": {
    "type": "object",
    "description": "Remove the specified prefixes from the URL path. It will strip the matching path prefix and will store the matching path prefix in a X-Forwarded-Prefix header.",
    "properties": {
      "prefixes": {
        "type": "array",
        "description": "The prefixes option defines the prefixes to strip from the request URL",
        "items": {
          "type": "string"
        }
      },
      "forceSlash": {
        "type": "boolean",
        "description": "The forceSlash option makes sure that the resulting stripped path is not the empty string, by replacing it with / when necessary.\n\nThis option was added to keep the initial (non-intuitive) behavior of this middleware, in order to avoid introducing a breaking change.\n\nIt's recommended to explicitly set forceSlash to false."
      }
    },
    "additionalProperties": false
  },
  "stripPrefixRegexMiddleware": {
    "type": "object",
    "description": "Remove the matching prefixes from the URL path. It will strip the matching path prefix and will store the matching path prefix in a X-Forwarded-Prefix header.",
    "properties": {
      "regex": {
        "type": "array",
        "description": "The regex option is the regular expression to match the path prefix from the request URL.",
        "items": {
          "type": "string"
        }
      }
    },
    "additionalProperties": false
  }
}

_.each(allSchemas, (schema) => {
  addPropertiesTitle(schema.properties)
})

export default allSchemas
