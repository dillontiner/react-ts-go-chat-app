import { createContext } from 'react'

const AuthContext = createContext({
    auth: "",
    setAuth: (v: string) => { },
})

export default AuthContext