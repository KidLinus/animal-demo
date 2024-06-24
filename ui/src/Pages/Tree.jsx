import { Flex, Spinner, Text } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import Input from "../Components/Input"
import AnimalRadialTree from "../Components/AnimalRadialTree"

const Tree = () => {
    const [state, stateSet] = useState({})
    const [animal, animalSet] = useState({})
    useEffect(() => {
        if (!state.animal) { animalSet({}); return }
        animalSet({ loading: true })
        fetch(`${import.meta.env.VITE_API_URI || "http://localhost:8666"}/animal/${state.animal}/family`).then(body => body.json())
            .then(data => data?.errors ? animalSet({ error: data.errors }) : animalSet({ data }))
            .catch(error => animalSet({ error }))
    }, [state.animal])
    return <Flex w="full" h="full" direction="column">
        <Flex gap="2" p="2" align="center" boxShadow="0 0 20px #0003">
            <Flex flexShrink="0">
                <Text>Tree page</Text>
            </Flex>
            <Input label="Animal A" value={state.animal || ""} onChange={v => stateSet({ ...state, animal: v })} />
        </Flex>
        <Flex w="full" h="full">
            {state.animal && <>
                {animal.loading && <Spinner />}
                {!animal.loading && animal.error && <Text>Loading animal A failed</Text>}
                {!animal.loading && !animal.error && <Flex w="full" h="full">
                    <AnimalRadialTree animal={animal.data} />
                </Flex>}
            </>}
        </Flex>
    </Flex>
}

export default Tree