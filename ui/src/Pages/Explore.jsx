import { Button, Flex, IconButton, Spacer, Spinner, Text } from '@chakra-ui/react'
import Input from "../Components/Input"
import Select from "../Components/Select"
import AnimalChart from "../Components/AnimalChart"
import { AiOutlinePlus, AiOutlineClose } from "react-icons/ai"
import { useEffect, useState } from 'react'

const Explore = () => {
  const [fields, fieldsSet] = useState({  })
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
      .then(data => data?.errors ? animalsSet({ error: data.errors }) : animalsSet({ data }))
      .catch(error => animalsSet({ error }))
  }, [fields])
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
        <Flex direction="column" gap="2">
          <Text>Animals</Text>
          {(fields?.animals || []).map((id,idx) => <Flex key={idx}>
            <Input value={id} onChange={n => fieldsSet({...fields, animals: fields.animals.map((v,i) => idx == i ? n : v)})} />
            <Spacer />
            <IconButton icon={<AiOutlineClose />} size="sm" variant="ghost" onClick={() => fieldsSet({...fields, animals: fields.animals.filter((_,i) => i !== idx)})} />
          </Flex>)}
          <Button leftIcon={<AiOutlinePlus />} onClick={() => fieldsSet({...fields, animals: [...(fields?.animals || []), 0]})}>Add</Button>
        </Flex>
      </>}
    </Flex>
    <Flex flexGrow="1" bg="gray.50">
      {(!fields.dataset || (fields?.animals || []).length < 1) 
        ? <Text>Select dataset and at least 1 animal</Text>
        : <>
        {animals.loading && <Spinner />}
        {!animals.loading && animals.error && <Text>Could not load animals</Text>}
        {!animals.loading && !animals.error && animals.data.length == 0 && <Text>No animals found</Text>}
        {!animals.loading && !animals.error && <AnimalChart data={animals.data} />}
      </>}
    </Flex>
  </Flex>
}

export default Explore
