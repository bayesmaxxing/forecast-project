import * as forecastService from './forecastService';
import * as pointsService from './pointsService';
import * as scoreService from './scoreService';
import * as userService from './userService';
//import * as blogService from './blogService';

// You can also define shared constants here
export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

export {
  forecastService,
  pointsService,
  scoreService,
  userService
 // blogService
};