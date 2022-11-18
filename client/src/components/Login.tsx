import React, { useState, useContext } from 'react';
import { styled } from '@mui/system';
import { Paper, TextField, Typography, Button, Link } from '@mui/material';
import Axios from 'axios'
import { useNavigate } from "react-router-dom"
import AuthContext from './AuthContext'

const LoginWindow = styled(Paper)({
    width: '20rem',
    height: '30rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
})

const LoginTitleContainer = styled('div')({
    paddingTop: '4rem',
    paddingBottom: '1rem',
})

const StyledLoginButton = styled(Button)({
    width: '10rem',
    textTransform: 'none',
})

const StyledForm = styled('form')({
    width: '100%',
    height: '100%',
    display: 'flex',
    justifyContent: 'space-evenly',
    alignItems: 'center',
    flexDirection: 'column',
    paddingTop: '2rem',
    paddingBottom: '2rem',
})

const SignUpLinkContainer = styled('div')({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
})

type LoginValues = {
    email: string
    password: string
}

const Login = () => {
    const navigate = useNavigate()
    const authContext = useContext(AuthContext)

    const [disabled, setDisabled] = useState(true)
    const [formValues, setFormValues] = useState<LoginValues>({
        email: '',
        password: '',
    })

    // TODO: consider adding auth to cookies

    const handleInputChange = (e: any) => {
        const { id, value } = e.target;
        const newFormValues = {
            ...formValues,
            [id]: value,
        }
        setFormValues(newFormValues)

        // TODO: better input validation
        if (newFormValues.email !== '' && newFormValues.password !== '') {
            setDisabled(false)
        } else {
            // TODO: handle delete case
            setDisabled(true)
        }
    };

    const handleSubmit = (event: any) => {
        event.preventDefault();

        // assumes input validation from handleInputChange
        const basicAuth = btoa(formValues.email + ":" + formValues.password) // TODO: hash password ?

        Axios({
            method: "GET",
            url: "http://127.0.0.1:4000/login",
            headers: {
                "Authorization": "Basic " + basicAuth
            }
        }).then(res => {
            // set auth context and go to the chat
            authContext.setAuth(res.data.uuid)
            navigate("/chat")
        }).catch((error) => {
            // TODO: handle errors
            console.log(error)
        })
    };

    return (
        <LoginWindow>
            <LoginTitleContainer>
                <Typography variant='h4'>Login</Typography>
            </LoginTitleContainer>
            <StyledForm onSubmit={handleSubmit}>
                <TextField
                    required
                    id='email'
                    label='email'
                    variant='filled'
                    onChange={handleInputChange}
                />
                <TextField
                    required
                    id='password'
                    label='password'
                    variant='filled'
                    onChange={handleInputChange}
                />
                <StyledLoginButton disabled={disabled} variant='contained' color='primary' type='submit'>
                    login
                </StyledLoginButton>
                <SignUpLinkContainer>
                    <Typography variant="body2">Don't have an account?</Typography>
                    <Typography variant="body2"><Link onClick={() => { navigate("/sign-up") }}>Sign Up</Link></Typography>
                </SignUpLinkContainer>
            </StyledForm>
        </LoginWindow>
    )
}

export default Login;
