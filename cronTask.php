<?php
/**消息格式示例*/
// $msgs = [
// 	[
// 		'date' => date('Y-m-d'),
// 		'source_id' => 'aaaaaaaaaaa',
// 		'source_type' => 'api',// 接口
// 		'urls' => [
// 			[
// 				'url' => 'https://weibo.com/ajax/statuses/mymblog',
// 				'params' => [
// 					'page' => '1',
// 					'feature' => '0',
// 					'uid' => '1270577247'
// 				],
// 				'fields' => [],
// 				'|' => [
// 					'url' => '',
// 					'params' => [
// 						'foo' => 'bar'
// 					],
// 					'fields' => [],
// 					'|' => [
// 						'url' => '',
// 						'params' => [
// 							'foo' => 'bar'
// 						],
// 						'fields' => [],
// 						'...'
//
// 					]
// 				]
// 			],
// 			[
// 				'...'
// 			],
// 		]
// 	],
// 	[
// 		'date' => date('Y-m-d'),
// 		'source_id' => 'bbbbbbbbbbb',
// 		'source_type' => 'file',// 文件
// 		'paths' => [
// 			[
// 				'account' => 'root',
// 				'pwd' => 'xxxxxxxxx', //decrypt
// 				'host' => 'xxxxxxxxxxxxx',
// 				'path' => "/root/test",
// 			 	'kind' => 'dir',
// 				'saveTo' => '/tmp/collect/multimedia'
// 			],
// 		]
// 	],
// 	[
// 		'date' => date('Y-m-d'),
// 		'source_id' => 'eeeeeeeee',
// 		'source_type' => 'csv',// csv
// 		'paths' => [
// 			[
// 				'account' => 'root',
// 				'pwd' => 'xxxxxxxxxxxxx', //decrypt
// 				'host' => 'xxxxxxxxxxxxxx',
// 				'path' => "/root/20220108.csv",
// 			 	'kind' => 'file',
// 				'saveTo' => '/tmp/collect/csv'
// 			],
// 		]
// 	],
// 	[
// 		'date' => date('Y-m-d'),
// 		'source_id' => 'cccccccccc',
// 		'source_type' => 'mongodb',// direct connect database
// 		'queries' => [
// 			[
// 				'query' => '',
// 				'params' => [
// 					'page' => '1',
// 					'feature' => '0',
// 					'uid' => '1270577247'
// 				],
// 			],
// 		]
// 	],
// 	[
// 		'date' => date('Y-m-d'),
// 		'source_id' => 'ddddddddddd',
// 		'source_type' => 'mysql',// direct connect database
// 		'queries' => [
// 			[
// 				'query' => '',
// 				'params' => [
// 					'page' => '1',
// 					'feature' => '0',
// 					'uid' => '1270577247'
// 				],
// 			],
// 		]
// 	],
//
// ];

$msgs = [];
$msgs[] = [
	'date' => date('Y-m-d'),
	'source_id' => 's323sasf23a46456qe2sdfd32x',
	'source_type' => 'csv',// multimedia file
	'paths' => [
		[
			'account' => 'root',
			'pwd' => 'xxxxxxxx', //decrypt
			'host' => 'tcp://xxxxxxx:22',
			'path' => "/root/testfile",
		 	'kind' => 'file',
			'saveTo' => '/tmp/collect/csv'
		],
		[
			'account' => 'root',
			'pwd' => 'xxxxxxx', //decrypt
			'host' => 'tcp://xxxxxx:22',
			'path' => '/root/test',
			'kind' => 'dir',

			'saveTo' => '/tmp/collect/multimedia'
		],
	]
];

$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
socket_connect($socket, '/tmp/collector.sock');

foreach($msgs as $msg) {
	$msg = json_encode($msg);
    socket_send($socket, $msg, strlen($msg), 0);
    $response = socket_read($socket, 1024);
	sleep(1);
    if($response != 'ok') {
		//todo log error
	}
}

socket_close($socket);






