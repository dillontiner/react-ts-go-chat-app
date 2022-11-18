

import React, { useState, useEffect, useContext } from 'react'
import { Paper, TextField, Button, Typography } from '@mui/material'
import ForwardIcon from '@mui/icons-material/Forward'
import ForwardOutlinedIcon from '@mui/icons-material/ForwardOutlined'
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
    marginTop: '0.5rem',
})

const DownVoted = styled(ForwardIcon)({
    transform: 'rotate(90deg)',
    color: 'red',
})

const UpVoted = styled(ForwardIcon)({
    transform: 'rotate(-90deg)',
    color: 'green',
})

const DownVoteEmpty = styled(ForwardOutlinedIcon)({
    transform: 'rotate(90deg)',
    color: 'blue',
})

const UpVoteEmpty = styled(ForwardOutlinedIcon)({
    transform: 'rotate(-90deg)',
    color: 'blue',
})

const NameVoteContainer = styled('div')({
    fontWeight: 'bold',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
})

type MessageProps = {
    message: Message
    sendVote: (vote: boolean, messageUuid: string) => void
}

type VoteProps = {
    votes: string[]
    sendMessageVote: (vote: boolean) => void
}

type Vote = {
    messageUuid: string
    upvoteUserUuids: string[]
    downvoteUserUuids: string[]
}

const VoteArrowContainer = styled('div')({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
})

const UpVote = ({ votes, sendMessageVote }: VoteProps) => {
    const n = votes.length
    const authContext = useContext(AuthContext)
    const userVoted = votes.includes(authContext.auth)
    return (
        <VoteArrowContainer>
            {userVoted ? (
                <UpVoted onClick={() => { sendMessageVote(true) }} />
            ) : (
                <UpVoteEmpty onClick={() => { sendMessageVote(true) }} />
            )}
            <>{n}</>
        </VoteArrowContainer>
    )
}

const DownVote = ({ votes, sendMessageVote }: VoteProps) => {
    const n = votes.length
    const authContext = useContext(AuthContext)
    const userVoted = votes.includes(authContext.auth)
    return (
        <VoteArrowContainer>
            {userVoted ? (
                <DownVoted onClick={() => { sendMessageVote(false) }} />
            ) : (
                <DownVoteEmpty onClick={() => { sendMessageVote(false) }} />
            )}
            <>{n}</>
        </VoteArrowContainer>
    )
}

const MessageDisplay = ({ message, sendVote }: MessageProps) => {
    // TODO: user upvoted or downvoted
    const sendMessageVote = (vote: boolean) => { sendVote(vote, message?.uuid || '') }
    return (
        <MessageContainer>
            <NameVoteContainer>
                {message.senderUuid}
                <NameVoteContainer>
                    <UpVote votes={message.upvoteUserUuids || []} sendMessageVote={sendMessageVote} />
                    <DownVote votes={message.downvoteUserUuids || []} sendMessageVote={sendMessageVote} />
                </NameVoteContainer>
            </NameVoteContainer>
            {message.body}
            <StyledTimestamp>{message.sentAt}</StyledTimestamp>
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
    const authContext = useContext(AuthContext)

    const sendVote = (vote: boolean, messageUuid: string) => {
        ws?.send(JSON.stringify({
            messageUuid: messageUuid,
            voterUuid: authContext.auth,
            vote: vote,
        }))
    }

    console.log(chatHistory)

    useEffect(() => {
        // TODO: query backend, redirect to login if failure
        if (ws != null) {
            ws.onmessage = function (evt: any) {
                const wsBody = JSON.parse(evt.data)?.body || {} // TODO: error handling
                const wsBodyJson = JSON.parse(wsBody)

                // MVP handling different types over one ws connection
                if (wsBodyJson["senderUuid"] !== undefined) {
                    setLastMessage(wsBodyJson as Message)
                    setLastMessageSentAt(new Date())
                } else if (wsBodyJson["voterUuid"] !== undefined) {
                    console.log(wsBodyJson)
                }


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
        if (lastMessage !== null) {
            setLiveMessages([lastMessage, ...liveMessages])
        }
    }, [lastMessageSentAt])

    return (
        <>
            {liveMessages.map((message) => (<MessageDisplay message={message} sendVote={sendVote} />))}
            {chatHistory.map((message) => (<MessageDisplay message={message} sendVote={sendVote} />))}
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

        var newWS: WebSocket | null = new WebSocket("ws://127.0.0.1:4001/chat")

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