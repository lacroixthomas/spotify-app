import { RootState } from '../../app/store';

import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface UserState {
  status: 'idle' | 'loading' | 'failed';
  username: string
  image: string
  id: string
}

const initialState: UserState = {
  status: 'idle',
  username: '',
  image: '',
  id: '',
};

export const getUserAsync = createAsyncThunk(
  'user/getUser',
  async (token: string) => {
    const response = await fetch('http://127.0.0.1:8080/user', { headers: { 'Authorization': token } })
    const json = await response.json();
    return json;
  }
);

export const authSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUser: (state, action: PayloadAction<object>) => ({
        ...state,
        ...action.payload,
    }),
  },
  extraReducers: (builder) => {
    builder
      .addCase(getUserAsync.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(getUserAsync.fulfilled, (state, action) => {
        state.status = 'idle';
        state.username = action.payload.name
        state.image = action.payload.image
        state.id = action.payload.id
      })
      .addCase(getUserAsync.rejected, (state, action) => {
        state.status = 'failed';
      });
  },
});

export const { setUser } = authSlice.actions;

export const selectUser = (state: RootState) => state.user;

export default authSlice.reducer;
