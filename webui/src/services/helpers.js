export function getProfileImage(fullPath) {
    const prefix = "webui/public";
    if (fullPath.startsWith(prefix)) {
        return fullPath.slice(prefix.length);
    }
    return fullPath;
}