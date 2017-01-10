package org.miaohong.jobs.msgjob.dal;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.hbase.HBaseConfiguration;
import org.apache.hadoop.hbase.HColumnDescriptor;
import org.apache.hadoop.hbase.HTableDescriptor;
import org.apache.hadoop.hbase.TableName;
import org.apache.hadoop.hbase.client.*;
import org.apache.hadoop.hbase.util.Bytes;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.concurrent.ThreadLocalRandom;

/**
 * Created by miaohong on 17/1/9.
 */
@Component
public class HBaseManager {
    private static final byte[] FAMILY = Bytes.toBytes("d");

    /**
     * For the example we're just using one qualifier.
     */
    private static final byte[] QUAL = Bytes.toBytes("test");

    private static Configuration configuration;
    private Connection connection;
    private String nameSpace;
    public HBaseManager() {
        configuration = HBaseConfiguration.create();
        configuration.set("hbase.master", "127.0.0.1:57336");
        configuration.set("hbase.zookeeper.property.clientPort", "2181");
        try {
            connection = ConnectionFactory.createConnection(configuration);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    public void createTable(String name){
        TableName tableName = TableName.valueOf(nameSpace, name);
        HTableDescriptor hTableDescriptor = new HTableDescriptor(tableName);
        HColumnDescriptor family = new HColumnDescriptor(FAMILY);
        hTableDescriptor.addFamily(family);
        try {
            connection.getAdmin().createTable(hTableDescriptor);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    public void insert(String tableName) {
        try {
//            HTable hTable = new HTable(configuration, tableName);
            byte[] rk = Bytes.toBytes(ThreadLocalRandom.current().nextLong());
            byte[] value = Bytes.toBytes(Double.toString(ThreadLocalRandom.current().nextDouble()));
            Put p = new Put(rk);
            p.addImmutable(FAMILY, QUAL, value);
            Table t = connection.getTable(TableName.valueOf(tableName));
            t.put(p);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
