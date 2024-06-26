import { Route, createBrowserRouter, createRoutesFromElements } from "react-router-dom"
import { Text } from "@chakra-ui/react"
import App from "./App"
import Browse from "./Pages/Browse"
import Animal from "./Pages/Animal"
import Child from "./Pages/Child"

const NotFound = () => <Text>Page not found</Text>

const router = createBrowserRouter(
    createRoutesFromElements(
        <Route path="/" element={<App />}>
            <Route index element={<Browse />} />
            <Route path="animal/:id" element={<Animal />} />
            <Route path="child/:a/:b" element={<Child />} />
            <Route path="*" element={<NotFound />} />
        </Route>
    ), { basename: import.meta.env.BASE_URL }
)

export default router