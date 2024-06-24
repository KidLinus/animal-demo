import { Flex, Spinner, Text } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import Input from "../Components/Input"
import AnimalTree from "../Components/AnimalTree"

const Tree3 = () => {
    const [state, stateSet] = useState({})
    const [animal1, animal1Set] = useState({})
    useEffect(() => {
        if (!state.animal1) { animal1Set({}); return }
        animal1Set({ loading: true })
        fetch(`${import.meta.env.VITE_API_URI || "http://localhost:8666"}/animal/${state.animal1}/family`).then(body => body.json())
            .then(data => data?.errors ? animal1Set({ error: data.errors }) : animal1Set({ data }))
            .catch(error => animal1Set({ error }))
    }, [state.animal1])
    const [animal2, animal2Set] = useState({})
    useEffect(() => {
        if (!state.animal2) { animal2Set({}); return }
        animal2Set({ loading: true })
        fetch(`${import.meta.env.VITE_API_URI || "http://localhost:8666"}/animal/${state.animal2}/family`).then(body => body.json())
            .then(data => data?.errors ? animal2Set({ error: data.errors }) : animal2Set({ data }))
            .catch(error => animal2Set({ error }))
    }, [state.animal2])
    return <Flex w="full" h="full" direction="column">
        <Flex gap="2" p="2" align="center" boxShadow="0 0 20px #0003">
            <Flex flexShrink="0">
                <Text>Tree page</Text>
            </Flex>
            <Input label="Animal A" value={state.animal1 || ""} onChange={v => stateSet({ ...state, animal1: v })} />
            <Input label="Animal B" value={state.animal2 || ""} onChange={v => stateSet({ ...state, animal2: v })} />
        </Flex>
        <Flex w="full" h="full">
            <Flex w="full" h="full" flexGrow="1">
                {state.animal1 && <>
                    {animal1.loading && <Spinner />}
                    {!animal1.loading && animal1.error && <Text>Loading animal1 A failed</Text>}
                    {!animal1.loading && !animal1.error && animal1.data && <Flex w="full" h="full">
                        <AnimalTree animal={animal1.data} />
                    </Flex>}
                </>}
            </Flex>
            <Flex w="full" h="full" flexGrow="1">
                {state.animal2 && <>
                    {animal2.loading && <Spinner />}
                    {!animal2.loading && animal2.error && <Text>Loading animal2 A failed</Text>}
                    {!animal2.loading && !animal2.error && animal2.data && <Flex w="full" h="full">
                        <AnimalTree animal={animal2.data} />
                    </Flex>}
                </>}
            </Flex>
        </Flex>
    </Flex>
}

export default Tree3