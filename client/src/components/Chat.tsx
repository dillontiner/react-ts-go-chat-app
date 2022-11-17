

import React, { useState, useEffect, useContext } from 'react'
import { Paper, TextField, Button, Typography } from '@mui/material'
import { styled } from '@mui/system'
import { useNavigate } from 'react-router'
import Axios from 'axios'
import AuthContext from './AuthContext'

const ChatWindow = styled(Paper)({
    width: '80%',
    height: '95%',
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
    height: '100%',
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

type Message = {
    body: string,
    sentAt: string,
    senderUuid: string,
    uuid?: string,
    upvoteUserUuids?: string[],
    downvoteUserUuids?: string[],
}

const ChatHistory = ({ ws }: ChatHistoryProps) => {
    const [lastMessage, setLastMessage] = useState<Message | null>(null)
    const [lastMessageSentAt, setLastMessageSentAt] = useState(new Date())
    const [chatHistory, setChatHistory] = useState<Message[]>([])
    const [liveMessages, setLiveMessages] = useState<Message[]>([])

    useEffect(() => {
        // TODO: query backend, redirect to login if failure
        if (ws != null) {
            ws.onmessage = function (evt: any) {
                setLastMessage(JSON.parse(evt.data) as Message)
                setLastMessageSentAt(new Date())
            }
        }
    }, [ws])

    useEffect(() => {
        Axios({
            method: "GET",
            url: "http://localhost:4000/chat",
        }).then(res => {
            // TODO: update messages in history
            setChatHistory(res?.data?.messages || [])
        }).catch((error) => {
            // TODO: handle errors
            console.log(error)
        })
    }, [])

    useEffect(() => {
        console.log(lastMessage)
        if (lastMessage !== null) {
            setLiveMessages([lastMessage, ...liveMessages])
        }
    }, [lastMessageSentAt])

    return (
        <>
            {liveMessages.map((message) => {
                console.log(message, message.body, message.senderUuid, message.sentAt)
                return (
                    <Message body={message.body} user={message.senderUuid} sentAt={message.sentAt} />
                )
            })}
            {chatHistory.map((message) => {
                return (
                    <Message body={message.body} user={message.senderUuid} sentAt={message.sentAt} />
                )
            })}
        </>
    )
}

const StyledForm = styled('form')({
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
})

const MessageTextField = styled(TextField)({
    width: '100%',
    height: '100%',
})
const SendButton = styled(Button)({
    width: '5rem',
    height: '2rem',
    textTransform: 'none',
    borderTopLeftRadius: '0px',
    borderTopRightRadius: '0px',
    alignSelf: 'flex-end',
})
const MessagePromptContainer = styled('div')({
    width: '100%',
    maxHeight: '10rem',
    marginTop: '1rem',
})

type MessagePromptProps = {
    ws: WebSocket | null
}
const MessagePrompt = ({ ws }: MessagePromptProps) => {
    const [disabled, setDisabled] = useState(true)
    const [messageBody, setMessageBody] = useState('')
    const authContext = useContext(AuthContext)

    const handleInputChange = (e: any) => {
        const newMessage = e.target.value
        setMessageBody(newMessage)

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

        // TODO: fix dates getting truncated (missing times, wrong order)
        const now = new Date()
        if (ws != null) {
            // TODO: handle writing on the backend
            const message: Message = {
                senderUuid: authContext.auth,
                body: messageBody,
                sentAt: now.toISOString()
            }
            ws.send(JSON.stringify(message))

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
                    maxRows={4}
                    id='messageBody'
                    variant='filled'
                    InputLabelProps={{ shrink: false }}
                    sx={{
                        borderTopRightRadius: '0px!important',
                        borderBottomRightRadius: '0px!important',
                    }}
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