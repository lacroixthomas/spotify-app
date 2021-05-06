import React from 'react';
import logo from './logo.svg';
import { Auth } from './features/auth/Auth';
import { Player } from './features/player/Player';
import './App.css';
import { useAppSelector } from './app/hooks';
import { selectToken } from './features/auth/authSlice';

function App() {

  const token = useAppSelector(selectToken);

  return (
    <div className="App">
      <div className="Left-panel">
        <div className="Login-container">

        TODO: To User Component -
        <br/>
        Image
        <br/>
        Username
        <br />
        <Auth />
        </div>
        <div className="Player-container">
          Player:

          Currently Playing ?

          Which music
          Which artist
          Which kind

          {token && <Player/>}          
        </div>
      </div>
      <div className="Playlist">
        Playlist:

        Listing of playlists from the user
      </div>
    </div>
  );
}

export default App;
