

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
    padding: '0.5rem',
    paddingLeft: '2rem',
    paddingRight: '2rem',
})

const ChatHistoryContainer = styled('div')({
    width: '100%',
    height: '20rem',
    display: 'flex',
    justifyContent: 'bottom',
    alignItems: 'left',
    flexDirection: 'column-reverse',
    overflowY: 'scroll',
})

const MessageContainer = styled('div')({
    padding: '1rem',
    borderBottom: '1px solid #CCCCCC',
    display: 'flex',
    flexDirection: 'column',
    overflow: 'wrap',
    wordBreak: 'break-all',
})

const StyledTimestamp = styled('div')({
    fontSize: '0.6rem',
})

type MessageProps = {
    body: string
    user: string
    sentAt: string
}

const Message = ({ body, user, sentAt }: MessageProps) => {
    return (
        <MessageContainer>
            <div><b>{user}</b></div>
            <div>{body}</div>
            <StyledTimestamp>{sentAt}</StyledTimestamp>
        </MessageContainer>
    )
}

type ChatHistoryProps = {
    ws: WebSocket | null
}


const ChatHistory = ({ ws }: ChatHistoryProps) => {
    const [lastMessage, setLastMessage] = useState("")
    const [lastMessageSentAt, setLastMessageSentAt] = useState(new Date())
    const [messages, setMessages] = useState<any[]>([])

    useEffect(() => {
        // TODO: query backend, redirect to login if failure
        if (ws != null) {
            ws.onmessage = function (evt: any) {
                setLastMessage(evt.data)
                setLastMessageSentAt(new Date())
            }
        }
    }, [ws])

    useEffect(() => {
        if (lastMessage !== "") {
            setMessages([lastMessage, ...messages])
        }
    }, [lastMessageSentAt])

    return (
        <>
            {messages.map((message) => <Message body={message} user={"TODO"} sentAt={"TODO"} />)}
        </>
    )
}

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
    textTransform: 'none',
})
const MessagePromptContainer = styled('div')({
    width: '100%',
})

type MessagePromptProps = {
    ws: WebSocket | null
}
const MessagePrompt = ({ ws }: MessagePromptProps) => {
    const [disabled, setDisabled] = useState(true)
    const [message, setMessage] = useState('')
    const authContext = useContext(AuthContext)

    const handleInputChange = (e: any) => {
        const newMessage = e.target.value
        setMessage(newMessage)

        if (newMessage !== '') {
            setDisabled(false)
        } else {
            setDisabled(true)
        }
    };

    if (ws == null) {
        // TODO: loading
        return (
            <>TODO: loading</>
        )
    }

    const handleSubmit = (event: any) => {
        event.preventDefault()

        // const now = new Date()
        if (ws != null) {
            // TODO: handle writing on the backend
            ws.send(message)

            // clear form inputs
            const form = document.getElementById("message-prompt") as HTMLFormElement
            if (form != null) {
                form.reset()
                setDisabled(true)
            }
        }
        // Axios({
        //     method: "POST",
        //     url: "http://localhost:4000/message",
        //     headers: {},
        //     data: {
        //         senderUuid: authContext.auth, // UUID to parametrize request
        //         body: message,
        //         sentAt: now.toISOString()
        //     }
        // }).then(res => {
        //     // TODO: update messages in app
        //     console.log(res)
        // }).catch((error) => {
        //     // TODO: handle errors
        //     console.log(error)
        // })
    };

    return (
        <MessagePromptContainer>
            <StyledForm id='message-prompt' onSubmit={handleSubmit}>
                <MessageTextField
                    multiline
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
    const [ws, setWS] = useState<WebSocket | null>(null)

    useEffect(() => {
        if (authContext.auth === "") {
            navigate("/login")
        }

        var newWS: WebSocket | null = new WebSocket("ws://127.0.0.1:4001/echo")

        // Improvement: better handling of ws connection
        newWS.onerror = function (evt) {
            console.log("ERROR: " + evt);
        }
        newWS.onclose = function (evt) {
            console.log("CLOSING");
            newWS = null
        }

        setWS(newWS)
    }, [])

    if (authContext.auth === "") {
        return <></> // redirecting
    }

    return (
        <ChatWindow>
            <Typography variant="h4">Chit Chat</Typography>
            <ChatHistoryContainer>
                <ChatHistory ws={ws} />
            </ChatHistoryContainer>
            <MessagePrompt ws={ws} />
        </ChatWindow>
    )
}

export default Chat