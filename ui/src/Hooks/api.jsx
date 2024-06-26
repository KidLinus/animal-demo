import { QueryClient, QueryClientProvider, onlineManager, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createContext, useContext } from "react"
import { stringify } from "qs"

export const apiCall = ({ method = "get", path = "/", body, query = {}, contentType = "json" }) => {
    const opts = { method: method.toUpperCase(), credentials: 'include' }
    if (body) {
        if (contentType === "formdata") {
            opts.body = toFormData(body)
        } else {
            opts.body = JSON.stringify(body)
            opts.headers = { "Content-Type": "application/json" }
        }
    }
    return fetch(`${import.meta.env.VITE_API_URI}${path}?${toQuery(query)}`, opts).then(res => {
        return Promise.resolve().then(() => {
            if ((res.headers.get("Content-Type") || "").startsWith("application/json")) { return res.json() }
            return res.text()
        }).then(v => {
            if (res.status == 200) { return v }
            if (typeof v == "object" && v?.id && v?.status) {
                throw new ApiError(v.id, v.status, v.details)
            }
            throw new ApiError("internal_server_error", v)
        })
    })
}

export class ApiError extends Error {
    constructor(id, status = 500, details = null) {
        super(`${id}${details ? " " + JSON.stringify(details) : ""}`)
        this.name = "ApiError"
        this.id = id
        this.status = status
        this.details = details
    }
}

export const toQuery = v => {
    const obj = Object.keys(v).reduce((s, k) => v[k] == "" ? s : { ...s, [k]: v[k] }, {})
    return stringify(obj, { encode: false, arrayFormat: "repeat", skipNulls: true })
}

export const toFormData = v => {
    const data = new FormData()
    Object.entries(v).forEach(([k, v]) => {
        if (Array.isArray(v) && v.length > 0 && v[0] instanceof File) {
            v.forEach(file => data.append(k, file))
            return
        }
        if (v instanceof File) { return data.append(k, v) }
        if (typeof v == "object" && !Array.isArray(v)) { return data.append(k, JSON.stringify(v)) }
        data.append(k, v)
    })
    return data
}

export const ApiContext = createContext()

export const useApi = () => useContext(ApiContext)

onlineManager.setOnline(true);
const queryClient = new QueryClient();

export const ApiProvider = ({ children }) => {
    return <ApiContext.Provider value={apiCall}>
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    </ApiContext.Provider>
}

const apiQueryKey = (query = {}) => {
    if (query?.queryKey) {return [query.queryKey]}
    const parts = query.path.split("/")
    if (parts.length >= 4) { return [parts[1], parts[2], parts[3], query.query || {}] }
    if (parts.length >= 3) { return [parts[1], parts[2], query.query || {}] }
    if (parts.length >= 2) { return [parts[1], query.query || {}] }
    return [query.query || {}]
}

export const useApiQuery = (query = {}, options = {}) => {
    return useQuery({ queryKey: apiQueryKey(query), queryFn: () => apiCall(query), ...options });
}

export const useApiMutation = (query = {}, options = {}) => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (opt) => apiCall({ ...query, ...opt }),
        onSuccess: (_, opt) => {
            queryClient.invalidateQueries(apiQueryKey({ ...query, ...opt }));
        },
        ...options,
    });
}

export const useInvalidate = () => {
    const queryClient = useQueryClient()
    return (path = "", query = {}) => {
        queryClient.invalidateQueries({ queryKey: apiQueryKey({ path, query }) })
    }
}
