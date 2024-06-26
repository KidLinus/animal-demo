import { AspectRatio, Button, Flex, Progress, Stat, StatHelpText, StatLabel, StatNumber, Table, Tbody, Td, Text, Th, Thead, Tr } from "@chakra-ui/react"
import { useNavigate, useParams } from "react-router-dom"
import { useApiQuery } from "../Hooks/api"
import { IoMdFemale, IoMdMale } from "react-icons/io"
import AnimalTree from "../Components/AnimalTree"
import { IoMaleFemale } from "react-icons/io5"

const Animal = () => {
    const { id } = useParams()
    const nav = useNavigate()
    const { data: animal, error, isPending } = useApiQuery({ path: `/animal/${id}` })
    const tree = useApiQuery({ path: `/animal/${id}/parents`, query: { depth: 6 } }, { enabled: !!animal })
    const coi = useApiQuery({ path: `/animal/${id}/coi`, query: { depth: 6 } }, { enabled: !!animal })
    return <Flex w="full" h="full" direction="column" p="2" gap="2">
        {isPending && <Progress isIndeterminate />}
        {!isPending && error && <Flex direction="column" gap="2">
            <Text>Could not load animal</Text>
            <Flex>
                <Button colorScheme="blue" onClick={() => nav("/")}>Go back</Button>
            </Flex>
        </Flex>}
        {!isPending && animal && <Flex direction="column" gap="2">
            <Flex gap="2" align="center">
                {animal.gender == "male" && <IoMdMale color="blue" />}
                {animal.gender == "female" && <IoMdFemale color="red" />}
                {animal.gender == "" && <IoMaleFemale color="#333" />}
                <Text>{animal.name}</Text>
            </Flex>
            <Text fontWeight="bold">Information</Text>
            <Flex direction="column">
                <Text>Born: {animal.born || "unknown"}</Text>
                <Text>Deceased: {animal.deceased || "unknown"}</Text>
            </Flex>
            <Text fontWeight="bold">Family tree</Text>
            <Flex>
                {tree?.data?.root ? <AspectRatio w="full" maxW="800px" ratio={4 / 3}>
                    <AnimalTree data={tree?.data?.root || {}} />
                </AspectRatio> : <Progress isIndeterminate />}
            </Flex>
            {coi?.data ? <Flex direction="column" gap="2">
                <Text fontWeight="bold">Inbreeding Coefficient</Text>
                <Stat>
                    <StatNumber>{Math.round(coi.data.result * 10000) / 100}%</StatNumber>
                    <StatHelpText>
                        Estimate
                    </StatHelpText>
                </Stat>
                <Text fontWeight="bold">Inbreeding Breakdown</Text>
                <Table size="sm" maxW="600px">
                    <Thead><Tr><Th>Parent</Th><Th>Relation</Th><Th>ICO</Th><Th>Contribution</Th></Tr></Thead>
                    <Tbody>
                        {coi.data.paths.map((path, idx) => <Tr key={idx}>
                            <Td>{path.parent}</Td>
                            <Td>{path.path.join("-")}</Td>
                            <Td>{Math.round(path.coi * 10000) / 100}%</Td>
                            <Td>{Math.round(path.result * 10000) / 100}%</Td>
                        </Tr>)}
                    </Tbody>
                </Table>
            </Flex> : <Progress isIndeterminate />}
            <Flex>
                <Button colorScheme="blue" onClick={() => nav("/")}>Go back</Button>
            </Flex>
        </Flex>
        }
    </Flex >
}

export default Animal