import React from 'react'
import ReactDOM from 'react-dom/client'
import { ChakraProvider } from '@chakra-ui/react'
import { RouterProvider } from "react-router-dom";
import { ApiProvider } from './Hooks/api';
import router from "./router"

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ApiProvider>
      <ChakraProvider>
        <RouterProvider router={router} />
      </ChakraProvider>
    </ApiProvider>
  </React.StrictMode>,
)
