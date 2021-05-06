import { RootState } from '../../app/store';

import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface PlaylistItem {
  image: string;
  name: string;
  ownerName: string;
  ID: string;
  URI: string;
}

export interface PlaylistState {
  status: 'idle' | 'loading' | 'failed';
  playlists: PlaylistItem[];
}

const initialState: PlaylistState = {
  status: 'idle',
  playlists: [],
};

export const getPlaylistAsync = createAsyncThunk(
  'playlist/getPlaylist',
  async (token: string) => {
    const response = await fetch('http://127.0.0.1:8080/playlist', { headers: { 'Authorization': token } })
    const json = await response.json();
    return json;
  }
);

export const startPlaylistAsync = createAsyncThunk(
  'playlist/startPlaylist',
  async (params: { token: string, uri: string }) => {

    const requestOptions = {
      method: 'POST',
      headers: {
        'Authorization': params.token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ uri: params.uri })
    };

    const response = await fetch('http://127.0.0.1:8080/player/play', requestOptions);
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
        if (!action.payload) {
          return;
        }
        const playlist = action.payload.map((item: any) => {
          return {
            image: item.image,
            name: item.name,
            ownerName: item.owner_name,
            ID: item.ID,
            URI: item.uri,
          } as PlaylistItem;
        });
        state.playlists = playlist;
      })
      .addCase(getPlaylistAsync.rejected, (state, action) => {
        state.status = 'failed';
      });
  },
});

export const { setPlaylist } = playlistSlice.actions;

export const selectPlaylist = (state: RootState) => state.playlist;

export default playlistSlice.reducer;
