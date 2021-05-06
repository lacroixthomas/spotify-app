import { useAppSelector, useAppDispatch } from '../../app/hooks';
import { selectToken } from '../auth/authSlice';
import { selectPlaylist, getPlaylistAsync, setPlaylist, PlaylistItem, startPlaylistAsync } from './playlistSlice';
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

  const items = playlist.playlists.map((item: PlaylistItem) => (
    <div className={styles.PlaylistItemContainer} key={item.ID}>
      <img className={styles.Image} src={item.image} onClick={() => dispatch(startPlaylistAsync({ token, uri: item.URI }))} ></img>
      <span className={styles.PlaylistItem} >{item.name}</span>
      <span className={styles.PlaylistItem} >by {item.ownerName}</span>
    </div>
  ));

  return (
    <div className={styles.PlaylistContainer}>
      {items}
    </div>
  );
}