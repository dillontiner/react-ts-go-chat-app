import React, { useState, useContext } from 'react';
import { styled } from '@mui/system';
import { Paper, TextField, Typography, Button } from '@mui/material';
import Axios from 'axios'
import { useNavigate } from "react-router-dom"
import AuthContext from './AuthContext'

const SignUpWindow = styled(Paper)({
    width: '20rem',
    height: '30rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
})

const SignUpTitleContainer = styled('div')({
    paddingTop: '4rem',
    paddingBottom: '1rem',
})

const StyledSignUpButton = styled(Button)({
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

type SignUpValues = {
    email: string
    password: string
    confirmPassword: string
}

const SignUp = () => {
    const navigate = useNavigate()
    const authContext = useContext(AuthContext)

    const [disabled, setDisabled] = useState(true)
    const [formValues, setFormValues] = useState<SignUpValues>({
        email: '',
        password: '',
        confirmPassword: '',
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
        if (newFormValues.email !== '' && newFormValues.password !== '' && newFormValues.password === formValues.confirmPassword) {
            setDisabled(false)
        } else {
            // TODO: handle delete case
            setDisabled(true)
        }
    }

    const handleSubmit = (event: any) => {
        event.preventDefault();

        // TODO: input validation

        // assumes input validation from handleInputChange
        const basicAuth = btoa(formValues.email + ":" + formValues.password) // TODO: hash password ?

        // TODO: hash passwords bcryptjs
        Axios({
            method: "POST",
            url: "http://localhost:4000/login",
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
        <SignUpWindow>
            <SignUpTitleContainer>
                <Typography variant='h4'>Sign Up</Typography>
            </SignUpTitleContainer>
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
                <TextField
                    required
                    id='confirmPassword'
                    label='confirm password'
                    variant='filled'
                    onChange={handleInputChange}
                />
                <StyledSignUpButton disabled={disabled} variant='contained' color='primary' type='submit'>
                    sign up
                </StyledSignUpButton>
            </StyledForm>
        </SignUpWindow>
    )
}

export default SignUp;
