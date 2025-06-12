import _ from 'lodash';
import axios from 'axios';

let baseUrl = _.trimEnd(location.pathname, '/')
export default axios.create({
  baseURL: `${baseUrl}/api/v1`,
})