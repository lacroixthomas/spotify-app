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
import ProgressBar from 'react-customizable-progressbar'
import playButtonPNG from './play-button.png';
import pauseButtonPNG from './pause-button.png';
import nextButtonPNG from './next-button.png';

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
      intervalID = (Number)(setInterval(() => dispatch(getPlayerAsync(token)), 1000));
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
      <div className={styles.musicInfo} >
        <p>{ player.isPlaying ? "Currently listening to:" : "No music playing"}</p>
        <p>Album: {player.albumName}</p>
        <p>Name: {player.musicName}</p>
        <p>Release date: {new Date(Date.parse(player.releaseDate)).toDateString()}</p>
        <p>Artists: {player.artists.map((a, ind) => (<span key={ind}>{a}<br /></span>))}</p>
      </div>
      <div className={styles.playerActionContainer}>
        <img className={styles.prevButton} onClick={() => dispatch(prevAsync(token))} src={nextButtonPNG}/>
        <ProgressBar
            className={styles.progressBar}
            progress={player.progress}
            radius={30}
            steps={player.duration}
            strokeColor="#1db954"
        >
          <img onClick={() => togglePlayPause(token)} className={styles.playButton} src={player.isPlaying ? pauseButtonPNG : playButtonPNG} />
        </ProgressBar>

        <img className={styles.nextButton} onClick={() => dispatch(nextAsync(token))} src={nextButtonPNG}/>
      </div>
    </div>
  );
}
