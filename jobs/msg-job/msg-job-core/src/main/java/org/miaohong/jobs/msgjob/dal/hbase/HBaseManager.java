package org.miaohong.jobs.msgjob.dal.hbase;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.hbase.HBaseConfiguration;
import org.apache.hadoop.hbase.HColumnDescriptor;
import org.apache.hadoop.hbase.HTableDescriptor;
import org.apache.hadoop.hbase.TableName;
import org.apache.hadoop.hbase.client.*;
import org.apache.hadoop.hbase.util.Bytes;
import org.miaohong.jobs.msgjob.dal.model.KafkaP2PMsg;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * Created by miaohong on 17/1/9.
 */
@Component
public class HBaseManager {
    @Value("${hbase.master}")
    private String hbaseMaster;
    @Value("${hbase.zookeeper.property.clientPort}")
    private String zkPort;

    @Value("${hbase.msg.table}")
    private String table;

    @Value("${hbase.msg.user.family}")
    private String userFamily;
    @Value("${hbase.msg.user.qual.sourceUID}")
    private String sourceUID;
    @Value("${hbase.msg.user.qual.targetUID}")
    private String targetUID;
    @Value("${hbase.msg.user.qual.online}")
    private String online;

    @Value("${hbase.msg.msg.family}")
    private String msgFamily;
    @Value("${hbase.msg.msg.qual.incrementID}")
    private String incrementID;
    @Value("${hbase.msg.msg.qual.msgID}")
    private String msgID;
    @Value("${hbase.msg.msg.qual.msg}")
    private String msg;
    @Value("${hbase.msg.msg.qual.accessServerAddr}")
    private String accessServerAddr;

    private static final Logger logger = LoggerFactory.getLogger(HBaseManager.class);


    //    private static final byte[] FAMILY = Bytes.toBytes("d");
//    private static final byte[] QUAL = Bytes.toBytes("test");
    private static Configuration configuration;
    private Connection connection;
    private String nameSpace;
    public void init() {
        configuration = HBaseConfiguration.create();
        configuration.set("hbase.master", hbaseMaster);
        configuration.set("hbase.zookeeper.property.clientPort", zkPort);
        try {
            connection = ConnectionFactory.createConnection(configuration);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
//    public void createTable(String name){
//        TableName tableName = TableName.valueOf(nameSpace, name);
//        HTableDescriptor hTableDescriptor = new HTableDescriptor(tableName);
//        HColumnDescriptor family = new HColumnDescriptor(FAMILY);
//        hTableDescriptor.addFamily(family);
//        try {
//            connection.getAdmin().createTable(hTableDescriptor);
//        } catch (IOException e) {
//            e.printStackTrace();
//        }
//    }
    public void insert(List<KafkaP2PMsg> kafkaP2PMsgs) {
        try {
            List<Put> puts = new ArrayList<>();
            Table t = connection.getTable(TableName.valueOf(table));
            for (KafkaP2PMsg kafkaP2PMsg : kafkaP2PMsgs) {
                byte[] rk = Bytes.toBytes(kafkaP2PMsg.getTargetUID() + "_" + kafkaP2PMsg.getIncrementID());
                // user
                byte[] sourceUIDData = Bytes.toBytes(kafkaP2PMsg.getSourceUID());
                byte[] incrementIDData = Bytes.toBytes(kafkaP2PMsg.getIncrementID());
                byte[] onlineData = Bytes.toBytes(kafkaP2PMsg.getOnline());
                byte[] accessServerAddrData = Bytes.toBytes(kafkaP2PMsg.getAccessServerAddr());

                // msg
                byte[] msgIDData = Bytes.toBytes(kafkaP2PMsg.getMsgID());
                byte[] msgData = Bytes.toBytes(kafkaP2PMsg.getMsg());
                Put p = new Put(rk);
                p.addImmutable(Bytes.toBytes(userFamily), Bytes.toBytes(sourceUID), sourceUIDData);
                p.addImmutable(Bytes.toBytes(userFamily), Bytes.toBytes(incrementID), incrementIDData);
                p.addImmutable(Bytes.toBytes(userFamily), Bytes.toBytes(online), onlineData);
                p.addImmutable(Bytes.toBytes(userFamily), Bytes.toBytes(accessServerAddr), accessServerAddrData);
                p.addImmutable(Bytes.toBytes(msgFamily), Bytes.toBytes(msgID), msgIDData);
                p.addImmutable(Bytes.toBytes(msgFamily), Bytes.toBytes(msg), msgData);
                puts.add(p);
            }
            t.put(puts);
        } catch (IOException e) {
            logger.error("insert err ", e);
        }
    }

    public void getUserMsg() {
        try {
            Table t = connection.getTable(TableName.valueOf(table));
//            t.get()
        } catch (IOException e) {
            logger.error("getUserMsg err ", e);
        }
    }
}
