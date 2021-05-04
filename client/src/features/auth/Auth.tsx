
import { useAppDispatch } from '../../app/hooks';
import { setToken } from './authSlice';

import styles from './Auth.module.css';
import { useEffect } from 'react';


const getHash = () => {
  const currentHash = window.location.hash
    .substring(1)
    .split("$")
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

export function Auth() {
  const dispatch = useAppDispatch();

  const endpoint = "https://accounts.spotify.com/authorize";
  const client_id = "e1cfaae8593c4e7b848e909c605e7ba3";
  const uri = "http://localhost:3000";

  useEffect(() => {
    const hash = getHash();
    if (!!hash['access_token']) {
      dispatch(setToken(hash['access_token']));
      // Store it in the local storage
    }
    // else check localstorage
  })

  return (
    <div>
        <a href={`${endpoint}?client_id=${client_id}&response_type=token&redirect_uri=${uri}`}>Login</a>
    </div>
  );
}
