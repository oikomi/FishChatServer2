package org.miaohong.jobs.msgjob.dal.model;

import com.alibaba.fastjson.annotation.JSONField;

/**
 * Created by miaohong on 17/1/14.
 */
public class KafkaP2PMsg {
    //@JSONField(name="sourceUID")
    private Long incrementID;
    private Long sourceUID;
    private Long targetUID;
    private String msgID;
    private String Msg;
    private String AccessServerAddr;
    private Boolean Online;

    public Long getIncrementID() {
        return incrementID;
    }

    public void setIncrementID(Long incrementID) {
        this.incrementID = incrementID;
    }

    public Long getSourceUID() {
        return sourceUID;
    }

    public void setSourceUID(Long sourceUID) {
        this.sourceUID = sourceUID;
    }

    public Long getTargetUID() {
        return targetUID;
    }

    public void setTargetUID(Long targetUID) {
        this.targetUID = targetUID;
    }

    public String getMsgID() {
        return msgID;
    }

    public void setMsgID(String msgID) {
        this.msgID = msgID;
    }

    public String getMsg() {
        return Msg;
    }

    public void setMsg(String msg) {
        Msg = msg;
    }

    public String getAccessServerAddr() {
        return AccessServerAddr;
    }

    public void setAccessServerAddr(String accessServerAddr) {
        AccessServerAddr = accessServerAddr;
    }

    public Boolean getOnline() {
        return Online;
    }

    public void setOnline(Boolean online) {
        Online = online;
    }

}
