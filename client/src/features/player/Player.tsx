import { useAppSelector, useAppDispatch } from '../../app/hooks';
import { selectToken } from '../auth/authSlice';
import {
  selectPlayer,
  getPlayerAsync,
  setPlayer,
  pauseAsync,
  playAsync,
  prevAsync,
  nextAsync,
  } from './playerSlice';
import styles from './Player.module.css';

import { useEffect } from 'react';

// intervalID interval ID of the refresh player
let intervalID = 0;

export function Player() {  
  const dispatch = useAppDispatch();
  const token = useAppSelector(selectToken);
  const player = useAppSelector(selectPlayer);

  useEffect(() => {
    if (!!token) {
      dispatch(getPlayerAsync(token));
      intervalID = (Number)(setInterval(() => dispatch(getPlayerAsync(token)), 5000));
    } else {
      dispatch(setPlayer({}));
      clearInterval(intervalID);
    }
    return () => {
      dispatch(setPlayer({}));
      clearInterval(intervalID);
    }
  }, [token, dispatch]);

  const togglePlayPause = (token: string) => {
    if (player.isPlaying) {
      dispatch(pauseAsync(token))
      dispatch(setPlayer({isPlaying: false}))
    } else {
      dispatch(playAsync(token))
      dispatch(setPlayer({isPlaying: true}))
    }
  }

  return (
    <div>
      { player.status === 'failed' && <span>An error occured</span> }
      <br />
      <span>{ player.isPlaying ? "Currently listening" : "No music playing"}</span>
      <br />
      {player.albumName}
      <br />
      {player.musicName}
      <br />
      Release date: {new Date(Date.parse(player.releaseDate)).toDateString()}
      <br />
      {player.artists.map((a, ind) => (<span key={ind}>{a}<br /></span>))}
      <br />

      ADD IMAGE OF CURRENT MUSIC

      <button onClick={() => dispatch(prevAsync(token))}>Previous</button>
      <button onClick={() => togglePlayPause(token)}>{player.isPlaying ? "Pause" : "Play"}</button>
      <button onClick={() => dispatch(nextAsync(token))}>Next</button>
    </div>
  );
}
