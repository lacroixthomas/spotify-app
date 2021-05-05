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
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <Auth />
        {token && <Player/>}
      </header>
    </div>
  );
}

export default App;
