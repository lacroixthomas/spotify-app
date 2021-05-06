import { useAppSelector, useAppDispatch } from '../../app/hooks';
import { selectToken } from '../auth/authSlice';
import { selectPlayer, getPlayerAsync, setPlayer } from './playerSlice';
import styles from './Player.module.css';

import { useEffect } from 'react';

export function Player() {  
  const dispatch = useAppDispatch();
  const token = useAppSelector(selectToken);
  const player = useAppSelector(selectPlayer);

  useEffect(() => {
    if (!!token) {
      dispatch(getPlayerAsync(token))
    } else {
      dispatch(setPlayer({}))
    }
  }, []);

  return (
    <div>
      { player.status == 'failed' && <span>An error occured</span> }
      { player.status == 'loading' && <span>Loading</span>}
      <br />
      <span>{ player.isPlaying ? "Currently listening" : "No music playing"}</span>
      <br />
      {player.albumName}
      <br />
      {player.musicName}
      <br />
      {player.artists.map((a, ind) => (<span key={ind}>{a}<br /></span>))}
      <br />
    </div>
  );
}
