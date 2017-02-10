package org.miaohong.jobs.msgjob.dal.model;

/**
 * Created by miaohong on 17/2/10.
 */
public class KafkaGroupMsg {
    private Long incrementID;
    private Long sourceUID;
    private Long targetUID;
    private Long groupID;
    private String msgID;
    private String Msg;

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

    public Long getGroupID() {
        return groupID;
    }

    public void setGroupID(Long groupID) {
        this.groupID = groupID;
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


//    private String AccessServerAddr;
//    private Boolean Online;

}
