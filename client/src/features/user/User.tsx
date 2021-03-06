import { useAppSelector, useAppDispatch } from '../../app/hooks';
import { selectToken } from '../auth/authSlice';
import { selectUser, getUserAsync, setUser } from './userSlice';
import styles from './User.module.css';

import { useEffect } from 'react';

export function User() {  
  const dispatch = useAppDispatch();
  const token = useAppSelector(selectToken);
  const user = useAppSelector(selectUser);

  useEffect(() => {
    if (!!token) {
      dispatch(getUserAsync(token))
    } else {
      dispatch(setUser({}))
    }
  }, [token, dispatch]);

  return (
    <div>
      { user.status === 'failed' && <span>An error occured</span> }
      { user.status === 'loading' && <span>Loading</span>}
      <img alt="user" className={styles.profile} src={user.image}></img>
      <p className={styles.username} >{user.username}</p>
    </div>
  );
}
