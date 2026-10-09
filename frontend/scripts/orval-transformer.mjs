// Orval input transformer: runs on backend/docs/swagger.json (Swagger 2.0, as
// swag emits it) before orval upgrades it to OpenAPI 3.
//
// - Drops the Go package prefix from schema names (`transport.FindUserResponse`
//   becomes `FindUserResponse`), so generated types read like the DTOs.
// - `transport.Response.data` is `omitempty` in Go, but swag's
//   --requiredByDefault still marks it required.

const PREFIX = "transport.";

const strip = (name) => (name.startsWith(PREFIX) ? name.slice(PREFIX.length) : name);

const rewriteRefs = (node) => {
  if (Array.isArray(node)) return node.map(rewriteRefs);
  if (node && typeof node === "object") {
    const out = {};
    for (const [key, value] of Object.entries(node)) {
      out[key] =
        key === "$ref" && typeof value === "string"
          ? value.replace(`#/definitions/${PREFIX}`, "#/definitions/")
          : rewriteRefs(value);
    }
    return out;
  }
  return node;
};

export default function transformer(spec) {
  const definitions = {};
  for (const [name, schema] of Object.entries(spec.definitions ?? {})) {
    definitions[strip(name)] = schema;
  }
  const result = rewriteRefs({ ...spec, definitions });

  const response = result.definitions.Response;
  if (response?.required) {
    response.required = response.required.filter((field) => field !== "data");
    if (response.required.length === 0) delete response.required;
  }

  return result;
}
