import { RootState } from '../../app/store';

import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface PlaylistState {
  status: 'idle' | 'loading' | 'failed';
}

const initialState: PlaylistState = {
  status: 'idle',
};

export const getPlaylistAsync = createAsyncThunk(
  'playlist/getPlaylist',
  async (token: string) => {
    const response = await fetch('http://127.0.0.1:8080/playlist', { headers: { 'Authorization': token } })
    const json = await response.json();
    return json;
  }
);

export const playlistSlice = createSlice({
  name: 'playlist',
  initialState,
  reducers: {
    setPlaylist: (state, action: PayloadAction<object>) => ({
        ...state,
        ...action.payload,
    }),
  },
  extraReducers: (builder) => {
    builder
      .addCase(getPlaylistAsync.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(getPlaylistAsync.fulfilled, (state, action) => {
        state.status = 'idle';
        console.log(action.payload)
      })
      .addCase(getPlaylistAsync.rejected, (state, action) => {
        state.status = 'failed';
      });
  },
});

export const { setPlaylist } = playlistSlice.actions;

export const selectPlaylist = (state: RootState) => state.playlist;

export default playlistSlice.reducer;
