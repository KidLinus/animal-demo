import { Flex } from "@chakra-ui/react"
import { Outlet } from "react-router-dom"

const App = () => {
  return <Flex w="full" h="full">
    <Outlet />
  </Flex>
}

export default App
