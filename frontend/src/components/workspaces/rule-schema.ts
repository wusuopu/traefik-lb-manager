import _ from 'lodash';

export const httpRouterSchema = {
  "type": "object",
  "description": "A router is in charge of connecting incoming requests to the services that can handle them. In the process, routers may use pieces of middleware to update the request, or act before forwarding the request to the service.",
  "properties": {
    "advanceMode": {
      "type": "boolean",
      "description": "If false, rule will just use `PathPrefix`."
    },
    "rule": {
      "type": "string",
      "description": "Rules are a set of matchers configured with values, that determine if a particular request matches specific criteria. If the rule is verified, the router becomes active, calls middlewares, and then forwards the request to the service."
    },
    "service": {
      "type": "integer",
      "description": "Each request must eventually be handled by a service, which is why each router definition should include a service target, which is basically where the request will be passed along to. HTTP routers can only target HTTP services (not TCP services).",
      "enum": [] as number[],
      "enumNames": [] as string[],
    },
    "priority": {
      "type": "integer",
      "description": "To avoid path overlap, routes are sorted, by default, in descending order using rules length. The priority is directly equal to the length of the rule, and so the longest length has the highest priority. A value of 0 for the priority is ignored: priority = 0 means that the default rules length sorting is used.",
      "default": 0,
      "minimum": 0,
      "maximum": 500
    },
    "middlewares": {
      "type": "array",
      "description": "You can attach a list of middlewares to each HTTP router. The middlewares will take effect only if the rule matches, and before forwarding the request to the service. Middlewares are applied in the same order as their declaration in router.",
      "items": {
        "enum": [] as number[],
        "enumNames": [] as string[],
        "type": "integer"
      },
      "uniqueItems": true,
    },
    "entryPoints": {
      "type": "array",
      "description": "If not specified, HTTP routers will accept requests from all defined entry points. If you want to limit the router scope to a set of entry points, set the entryPoints option.",
      "items": {
        "type": "string",
        "enum": [] as string[],
      },
      "uniqueItems": true,
    },
  },
  "additionalProperties": false,
  "required": [
    "rule",
    "service"
  ]
}
function addPropertiesTitle (properties: any) {
  for (const key in properties) {
    const item = properties[key]
    if (!item.title) { item.title = key }

    if (item.properties) {
      addPropertiesTitle(item.properties)
    }
  }
}
addPropertiesTitle(httpRouterSchema.properties)