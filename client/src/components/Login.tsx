import React, { useState } from 'react';
import { styled } from '@mui/system';
import { Paper, TextField, Typography, Button } from '@mui/material';

const LoginWindow = styled(Paper)({
    width: '20rem',
    height: '30rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
})

const LoginTitleContainer = styled('div')({
    padding: '4rem',
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

type LoginValues = {
    email: string
    password: string
}

const Login = () => {
    const [formValues, setFormValues] = useState<LoginValues>({
        email: '',
        password: '',
    })

    const handleInputChange = (e: any) => {
        const { id, value } = e.target;
        setFormValues({
            ...formValues,
            [id]: value,
        });
    };

    const handleSubmit = (event: any) => {
        event.preventDefault();
        console.log(formValues);
        // TODO: input validation
        if (formValues.email != '' && formValues.password != '') {
            console.log('should be good')
        }
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
                <StyledLoginButton variant='contained' color='primary' type='submit'>
                    login
                </StyledLoginButton>
            </StyledForm>
        </LoginWindow>
    )
}

export default Login;
