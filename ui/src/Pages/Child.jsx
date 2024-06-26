import { AspectRatio, Button, Flex, Progress, Stat, StatHelpText, StatLabel, StatNumber, Table, Tbody, Td, Text, Th, Thead, Tr } from "@chakra-ui/react"
import { useNavigate, useParams } from "react-router-dom"
import { useApiQuery } from "../Hooks/api"
import { IoMdFemale, IoMdMale } from "react-icons/io"
import { IoMaleFemale } from "react-icons/io5"
import AnimalTree from "../Components/AnimalTree"
import { useState } from "react"
import Input from "../Components/Input"

const Child = () => {
    const { a, b } = useParams()
    const nav = useNavigate()
    const [depth, depthSet] = useState("6")
    const animal = { id: "temporary", name: "Potential Child", gender: "", parents: { male: a, female: b } }
    const tree = useApiQuery({ path: `/animal/parents`, query: { a, b, depth: Math.min(parseInt(depth, 10), 14) } }, { enabled: !!animal && parseInt(depth, 10) > 0 })
    const coi = useApiQuery({ path: `/animal/coi`, query: { a, b, depth: Math.min(parseInt(depth, 10), 14) } }, { enabled: !!animal && parseInt(depth, 10) > 0 })
    return <Flex w="full" h="full" direction="column" p="2" gap="2">
        <Flex direction="column" gap="2">
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
            <Text fontWeight="bold">Analysis Depth (max 14)</Text>
            <Input label="Depth" type="number" min="1" value={depth} onChange={depthSet} maxW="150px" />
            <Flex>
                <Button colorScheme="blue" onClick={() => nav("/")}>Go back</Button>
            </Flex>
        </Flex>
    </Flex >
}

export default Child