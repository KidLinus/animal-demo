import { Button, Flex, Progress, Link, Table, Tbody, Td, Text, Th, Thead, Tr } from "@chakra-ui/react"
import { Link as RouterLink } from "react-router-dom"
import { useState } from "react"
import { IoMdFemale, IoMdMale } from "react-icons/io"
import Input from "../Components/Input"
import { useApiMutation } from "../Hooks/api"

const Browse = () => {
    const [state, stateSet] = useState({})
    const { mutateAsync: search, data, isPending } = useApiMutation({ path: "/animal", query: state })
    return <Flex w="full" h="full" direction="column" p="2" gap="2">
        <Flex gap="2" as="form" onSubmit={e => { e.preventDefault(); search() }}>
            <Input label="Search animal" value={state.query} onChange={query => stateSet({ ...state, query })} />
            <Button type="submit" colorScheme="blue" isLoading={isPending}>Search</Button>
        </Flex>
        <Flex gap="2" direction="column">
            {isPending && <Progress isIndeterminate />}
            {!isPending && data?.items?.length == 0 && <Text>No results found</Text>}
            {!isPending && data?.items?.length > 0 && <Table>
                <Thead><Tr><Th>Name</Th><Th>Gender</Th><Th>Born</Th><Th>Deceased</Th></Tr></Thead>
                <Tbody>
                    {data.items.map(animal => <Tr key={animal.id}>
                        <Td><Link to={`/animal/${animal.id}`} as={RouterLink}>{animal.name}</Link></Td>
                        <Td>{animal.gender == "male" && <IoMdMale color="blue" />}{animal.gender == "female" && <IoMdFemale color="red" />}</Td>
                        <Td>{animal.born || "-"}</Td>
                        <Td>{animal.deceased || "-"}</Td>
                    </Tr>)}
                </Tbody>
            </Table>}
        </Flex>
    </Flex>
}

export default Browse