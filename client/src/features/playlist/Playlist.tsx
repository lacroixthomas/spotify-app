import { useAppSelector, useAppDispatch } from '../../app/hooks';
import { selectToken } from '../auth/authSlice';
import { selectPlaylist, getPlaylistAsync, setPlaylist } from './playlistSlice';
import styles from './Playlist.module.css';

import { useEffect } from 'react';

export function Playlist() {  
  const dispatch = useAppDispatch();
  const token = useAppSelector(selectToken);
  const playlist = useAppSelector(selectPlaylist);

  useEffect(() => {
    if (!!token) {
      dispatch(getPlaylistAsync(token))
    } else {
      dispatch(setPlaylist({}))
    }
  }, []);

  return (
    <div>
      {JSON.stringify(playlist)}
      <br />
    </div>
  );
}
