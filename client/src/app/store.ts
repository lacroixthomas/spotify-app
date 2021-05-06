import { configureStore, ThunkAction, Action } from '@reduxjs/toolkit';
import authReducer from '../features/auth/authSlice';
import userReducer from '../features/user/userSlice';
import playerReducer from '../features/player/playerSlice';
import playlistReducer from '../features/playlist/playlistSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    user: userReducer,
    player: playerReducer,
    playlist: playlistReducer,
  },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
