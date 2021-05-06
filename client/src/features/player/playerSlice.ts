import { RootState } from '../../app/store';

import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface PlayerState {
  status: 'idle' | 'loading' | 'failed';
  isPlaying: boolean;
  albumName: string;
  musicName: string
  artists: string[];
  id: string;
}

const initialState: PlayerState = {
  status: 'idle',
  isPlaying: false,
  musicName: '',
  albumName: '',
  artists: [],
  id: '',
};

export const getPlayerAsync = createAsyncThunk(
  'player/getPlayer',
  async (token: string) => {
    const response = await fetch('http://127.0.0.1:8080/player', { headers: { 'Authorization': token } })
    const json = await response.json();
    return json;
  }
);

export const playerSlice = createSlice({
  name: 'player',
  initialState,
  reducers: {
    setPlayer: (state, action: PayloadAction<object>) => ({
        ...state,
        ...action.payload,
    }),
  },
  extraReducers: (builder) => {
    builder
      .addCase(getPlayerAsync.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(getPlayerAsync.fulfilled, (state, action) => {
        state.status = 'idle';
        state.albumName = action.payload.album_name;
        state.artists = action.payload.artists_name || [];
        state.musicName = action.payload.music_name;
        state.isPlaying = action.payload.is_playing;
        state.id = action.payload.id;
      })
      .addCase(getPlayerAsync.rejected, (state, action) => {
        state.status = 'failed';
      });
  },
});

export const { setPlayer } = playerSlice.actions;

export const selectPlayer = (state: RootState) => state.player;

export default playerSlice.reducer;
