#!/usr/bin/env bash

###############################################################################
# Check for USB errors and restart services if needed
###############################################################################

echo "Starting $0"

# Source common utilities
# shellcheck source=../utils.sh
source "$(dirname "${BASH_SOURCE[0]}")/../utils.sh"

# Log also to a file
LogFile="/tmp/check_usb_errors/check_usb_errors.log"
exec > >(stdbuf -o0 sed 's/%/%%/g' | xargs -d '\n' -I {} date '+%F %T {}' | tee -a "$LogFile" "${LogFile}-$(date --iso-8601=seconds)")
exec 2>&1

# TODO: clean old log files every once in a while

restart_deployments() {
    log_info "Restarting Jellyfin and qBittorrent deployments in the 'server' namespace..."
    kubectl rollout restart "deployment/jellyfin" -n server
    kubectl rollout restart "deployment/qbittorrent" -n server
    log_info "Deployments accessing media storage were restarted"
}

main() {
    log_info "Checking for USB errors..."
    
    log_info "Looking through 'dmesg' logs..."
    usb_logs=$(dmesg -T | grep "usb" | tail)
    
    if [[ -z "$usb_logs" ]]; then
        log_info "No USB events detected"
        return 0
    fi
    
    log_info "Checking for USB reconnect events..."
    reconnect_events=$(echo "$usb_logs" | grep "new high-speed USB device number")
    
    if [[ -z "$reconnect_events" ]]; then
        log_info "No USB reconnect events detected"
        return 0
    fi
    
    log_info "Detected USB disconnect and reconnect events"
    
    log_info "Parsing last reconnect time..."
    reconnect_time=$(echo "$reconnect_events" | tail -1 | grep -o '\[.*\]' | tr -d '[]')
    
    log_info "Converting the dmesg timestamp to Unix timestamp..."
    usb_unix_time=$(date -d "$reconnect_time" +%s 2>/dev/null)
    
    # Exit if we couldn't parse the timestamp
    if [[ -z "$usb_unix_time" ]]; then
        log_error "Failed to parse reconnect timestamp"
        return 1
    fi
    
    log_info "Last USB device reconnected at: $reconnect_time (Unix timestamp: $usb_unix_time)"
    
    log_info "Checking deployment timestamps..."
    
    jellyfin_k8s_time=$(kubectl get deployment jellyfin -n server -o jsonpath='{.metadata.creationTimestamp}')
    jellyfin_unix_time=$(date -d "$jellyfin_k8s_time" +%s 2>/dev/null)
    log_info "Jellyfin deployment created at: $jellyfin_k8s_time (Unix timestamp: $jellyfin_unix_time)"
    
    usb_unix_time_int=$((usb_unix_time))  # Convert string to integer
    jellyfin_unix_time_int=$((jellyfin_unix_time))  # Convert string to integer
    
    if [[ $jellyfin_unix_time_int -ge $usb_unix_time_int ]]; then
        log_info "Deployments are newer than USB reconnect event. No restart needed."
        return 0
    fi
    
    log_info "USB reconnect event is newer than deployments. Restarting affected deployments..."
    restart_deployments

    return 0
}

main "$@"

echo "Finished $0"
