import _ from 'lodash';
import moment from 'moment';

const datetimeFormat = (date: string, format: string = 'YYYY-MM-DD HH:mm:ss'): string => {
  if (!date) { return '' }
  if (_.startsWith(date, '0001-')) { return '' }
  return moment(date).format(format);
}

export default {
  datetimeFormat,
  tableDatetimeFormat: (row: any, column: any, cellValue: any, index: any): string => {
    return datetimeFormat(cellValue)
  }
}