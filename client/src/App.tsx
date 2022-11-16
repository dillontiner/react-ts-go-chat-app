import React, { useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import {
  Routes,
  Route,
  Navigate
} from 'react-router';
import { styled } from '@mui/system';
import Login from './components/Login'

const StyledApp = styled('div')({
  background: 'black',
  width: '100vw',
  height: '100vh',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
})

const TODO = () => {
  return (
    <div>TODO</div>
  )
}

const App = () => {
  return (
    <StyledApp>
      <Router>
        <Routes>
          <Route path='/login' element={<Login />} />
          <Route path='/sign-up' element={<TODO />} />
          <Route path='/chat' element={<TODO />} />
          <Route path='*' element={<Navigate to='/login' />} />
        </Routes>
      </Router >
    </StyledApp>
  );
}

export default App;
