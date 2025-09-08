package main

import "github.com/gin-gonic/gin"

type Permission int

const (
    PermissionRead Permission = 1 << iota
    PermissionWrite
    PermissionAdmin
)

var apiKeys = map[string]Permission{
    "admin-key":  PermissionRead | PermissionWrite | PermissionAdmin, // 全部权限
    "writer-key": PermissionRead | PermissionWrite,                   // 读写权限
    "reader-key": PermissionRead,                                    // 只读权限
}

func validateAPIKey(key string, requiredPerm Permission) bool {
    if perm, exists := apiKeys[key]; exists {
        return perm&requiredPerm != 0
    }
    return false
}

func getAPIKeyFromHeader(c *gin.Context) string {
    return c.GetHeader("X-API-Key")
}
