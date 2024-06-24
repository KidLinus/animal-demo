import { Route, createBrowserRouter, createRoutesFromElements } from "react-router-dom"
import { Text } from "@chakra-ui/react"
import App from "./App"
import Explore from "./Pages/Explore"
import Tree from "./Pages/Tree"
import Tree2 from "./Pages/Tree2"

const NotFound = () => <Text>Page not found</Text>

const router = createBrowserRouter(
    createRoutesFromElements(
        <Route path="/" element={<App />}>
            <Route index element={<Explore />} />
            <Route path="tree" element={<Tree />} />
            <Route path="tree2" element={<Tree2 />} />
            <Route path="*" element={<NotFound />} />
        </Route>
    ), { basename: import.meta.env.BASE_URL }
)

export default router