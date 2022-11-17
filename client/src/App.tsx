import React, { useContext, useEffect, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import {
  Routes,
  Route,
  Navigate,
  useNavigate
} from 'react-router';
import { styled } from '@mui/system';
import Login from './components/Login'
import SignUp from './components/SignUp'
import AuthContext from './components/AuthContext'
import Chat from './components/Chat'

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
  const [auth, setAuth] = useState("")

  return (
    <AuthContext.Provider value={{ auth: auth, setAuth: setAuth }}>
      <StyledApp>
        <Router>
          <Routes>
            <Route path='/login' element={<Login />} />
            <Route path='/sign-up' element={<SignUp />} />
            <Route path='/chat' element={<Chat />} />
            <Route path='*' element={<Navigate to='/login' />} />
          </Routes>
        </Router >
      </StyledApp>
    </AuthContext.Provider >
  );
}

export default App;
