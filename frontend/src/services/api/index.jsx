import * as forecastService from './forecastService';
import * as pointsService from './pointsService';
import * as scoreService from './scoreService';
import * as userService from './userService';

// You can also define shared constants here
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export {
  forecastService,
  pointsService,
  scoreService,
  userService
};