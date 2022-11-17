

import React, { useState, useEffect, useContext } from 'react'
import { Paper, TextField, Button, Typography } from '@mui/material'
import { styled } from '@mui/system'
import { useNavigate } from 'react-router'
import Axios from 'axios'
import AuthContext from './AuthContext'

const ChatWindow = styled(Paper)({
    width: '30rem',
    height: '30rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
    margin: '0.5rem',
    marginLeft: '2rem',
    marginRight: '2rem',
})

const ChatHistory = styled('div')({
    width: '30rem',
    height: '20rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
})

const StyledForm = styled('form')({
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-evenly',
})

const MessageTextField = styled(TextField)({
    width: '100',
})
const SendButton = styled(Button)({
    width: '2rem',
})
const MessagePromptContainer = styled('div')({
    width: '100%',
})

const MessagePrompt = () => {
    const [disabled, setDisabled] = useState(true)
    const [message, setMessage] = useState('')
    const authContext = useContext(AuthContext)

    // TODO: consider adding auth to cookies

    const handleInputChange = (e: any) => {
        const newMessage = e.target.value
        setMessage(newMessage)

        if (newMessage !== '') {
            setDisabled(false)
        } else {
            setDisabled(true)
        }
    };

    const handleSubmit = (event: any) => {
        event.preventDefault()

        const now = new Date()
        console.log(now)
        Axios({
            method: "POST",
            url: "http://localhost:4000/message",
            data: {
                senderUuid: authContext.auth, // UUID to parametrize request
                message: message,
                sentAt: now.toISOString()
            }
        }).then(res => {
            // set auth context and go to the chat
            console.log(res)
        }).catch((error) => {
            // TODO: handle errors
            console.log(error)
        })
    };

    return (
        <MessagePromptContainer>
            <StyledForm onSubmit={handleSubmit}>
                <MessageTextField
                    multiline
                    rows={2}
                    maxRows={10}
                    id='message'
                    variant='filled'
                    onChange={handleInputChange}
                />
                <SendButton disabled={disabled} variant='contained' color='primary' type='submit'>
                    send
                </SendButton>
            </StyledForm>
        </MessagePromptContainer>
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
        <ChatWindow>
            <Typography variant="h4">Chit Chat</Typography>
            <ChatHistory />
            <MessagePrompt />
        </ChatWindow>
    )
}

export default Chat