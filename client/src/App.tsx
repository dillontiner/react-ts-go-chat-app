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
import AuthContext from './components/AuthContext'

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


const Chat = () => {
  const authContext = useContext(AuthContext)
  const navigate = useNavigate()

  useEffect(() => {
    // TODO: query backend, redirect to login if failure
    if (authContext.auth === "") {
      navigate("/login")
    }
  })

  return (
    <TODO />
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
            <Route path='/sign-up' element={<TODO />} />
            <Route path='/chat' element={<Chat />} />
            <Route path='*' element={<Navigate to='/login' />} />
          </Routes>
        </Router >
      </StyledApp>
    </AuthContext.Provider >
  );
}

export default App;
