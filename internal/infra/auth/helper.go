package auth


func getStringClaim(claims map[string]interface{}, key string) string {
    if val, ok := claims[key]; ok {
        if str, ok := val.(string); ok {
            return str
        }
    }
    return ""
}
