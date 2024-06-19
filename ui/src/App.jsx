import { Button, Flex, IconButton, Spacer, Text } from '@chakra-ui/react'
import Select from "./Components/Select"
import AnimalChart from "./Components/AnimalChart"
import { AiOutlinePlus, AiOutlineClose } from "react-icons/ai"
import { useEffect, useState } from 'react'

const App = () => {
  const [fields, fieldsSet] = useState({ dataset: 1, animals: [526094, 428660] })
  const [datasets, datasetsSet] = useState({ loading: true })
  useEffect(() => {
    fetch(`${import.meta.env.VITE_API_URI || "http://localhost:8666"}/dataset`).then(body => body.json())
      .then(data => datasetsSet({ data }))
      .catch(error => datasetsSet({ error }))
  }, [])
  const [animals, animalsSet] = useState({ loading: true })
  useEffect(() => {
    if (!fields.dataset || (fields?.animals || []).length < 1) {return}
    fetch(`${import.meta.env.VITE_API_URI || "http://localhost:8666"}/animal/family?${(fields?.animals || []).map(v => `id=${v}`).join("&")}`).then(body => body.json())
      .then(data => animalsSet({ data }))
      .catch(error => animalsSet({ error }))
  }, [fields])
  console.log(fields, animals)
  return <Flex w="full" h="full">
    <Flex w="200px" zIndex="1" boxShadow="0 0 20px #0003" p="2" direction="column" gap="4">
      <Flex direction="column">
        <Text>Dataset</Text>
        {!!datasets.loading && <Text fontSize="sm">Loading...</Text>}
        {!datasets.loading && !!datasets.error && <Text fontSize="sm">Failed to load data sets</Text>}
        {!datasets.loading && !datasets.error && <Select value={fields.dataset == null ? "": fields.dataset}
          onChange={v => fieldsSet({ dataset: v == "" ? null : parseInt(v, 10) })}>
          <option value="">None selected</option>
          {datasets.data.map(set => <option key={set.id} value={set.id}>{set.name}</option>)}
        </Select>}
      </Flex>
      {!!fields.dataset && <>
        <Flex direction="column">
          <Text>Animals</Text>
          {(fields?.animals || []).map(id => <Flex key={id}>
            <Text>{id}</Text>
            <Spacer />
            <IconButton icon={<AiOutlineClose />} size="sm" variant="ghost" />
          </Flex>)}
          <Button leftIcon={<AiOutlinePlus />}>Add</Button>
        </Flex>
      </>}
    </Flex>
    <Flex flexGrow="1" bg="gray.50">
      <AnimalChart data={(animals?.data || [])} />
    </Flex>
  </Flex>
}

export default App
