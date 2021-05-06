import React from 'react';
import logo from './logo.svg';
import { Auth } from './features/auth/Auth';
import { Player } from './features/player/Player';
import { Playlist } from './features/playlist/Playlist';
import { User } from './features/user/User';
import './App.css';
import { useAppSelector } from './app/hooks';
import { selectToken } from './features/auth/authSlice';

function App() {

  const token = useAppSelector(selectToken);

  return (
    <div className="App">
      <div className="Left-panel">
        <div className="Login-container">
        {!! token && <User /> }
        <Auth />
        </div>
        <div className="Player-container">
          {token && <Player/>}
        </div>
      </div>
      <div className="Playlist">
      {token && <Playlist/>}
      </div>
    </div>
  );
}

export default App;
