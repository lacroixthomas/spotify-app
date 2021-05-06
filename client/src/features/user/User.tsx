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
  }, []);

  return (
    <div>
      { user.status == 'failed' && <span>An error occured</span> }
      { user.status == 'loading' && <span>Loading</span>}
      <img className={styles.profile} src={user.image}></img>
      <br />
      <span>{user.username}</span>
    </div>
  );
}
