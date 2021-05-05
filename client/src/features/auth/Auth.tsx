
import { useAppDispatch } from '../../app/hooks';
import { setToken, selectToken } from './authSlice';
import { useAppSelector } from '../../app/hooks';

import styles from './Auth.module.css';
import { useEffect } from 'react';

export function Auth() {
  const dispatch = useAppDispatch();
  const token = useAppSelector(selectToken);
  const endpoint = "https://accounts.spotify.com/authorize";
  const client_id = "e1cfaae8593c4e7b848e909c605e7ba3";
  const uri = "http://localhost:3000";

  // Retrieve token if possible (hash or local storage)
  useEffect(() => {
    let accessToken = window.localStorage.getItem('spotify-access-token');
    const hash = getHash();
    if (!!hash['access_token']) {
      accessToken = hash['access_token'];
    }
    if (accessToken) {
      dispatch(setToken(accessToken));
      window.localStorage.setItem('spotify-access-token', accessToken);
      // TODO: Set refresh token every x seconds
    } else {
      logOut()
    }
  })

  // getHash parse the hash and return a map of all hash
  const getHash = () => {
    const currentHash = window.location.hash
      .substring(1)
      .split("&")
      .reduce((acc: any, next) => {
        if (next) {
          const hashParts = next.split("=");
          acc[hashParts[0]] = decodeURIComponent(hashParts[1]);
        }
        return acc;
      }, {});
    window.location.hash = '';
    return currentHash;
  }

  // logOut will log the current user out (cleaning token)
  const logOut = () => {
    window.localStorage.removeItem('spotify-access-token')
    dispatch(setToken(""));
  }

  // scopes to use when connecting to spotify api
  const scopes = [
    "playlist-read-private",
    "playlist-read-collaborative",
    "playlist-modify-public",
    "user-read-recently-played",
    "playlist-modify-private",
    "ugc-image-upload",
    "user-follow-modify",
    "user-follow-read",
    "user-library-read",
    "user-library-modify",
    "user-read-private",
    "user-read-email",
    "user-top-read",
    "user-read-playback-state"
  ];

  return (
    <div>
        { !token &&
          <a href={`${endpoint}?client_id=${client_id}&scope=${encodeURIComponent(scopes.join(' '))}&response_type=token&redirect_uri=${uri}`}>Login</a>
        }
        { token &&
          <button onClick={logOut}>Log out</button>
        }
    </div>
  );
}
